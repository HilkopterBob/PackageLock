package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/server"
	"packagelock/structs"
	"strconv"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	restartChan = make(chan struct{})
	quitChan    = make(chan os.Signal, 1)
	AppVersion  string // Version injected with ldflags
)

// Root command using Cobra
var rootCmd = &cobra.Command{
	Use:   "packagelock",
	Short: "Packagelock CLI tool",
	Long:  `Packagelock CLI manages the server and other operations.`,
}

// Start command to run the server
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the running server",
	Run: func(cmd *cobra.Command, args []string) {
		restartServer()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running server",
	Run: func(cmd *cobra.Command, args []string) {
		stopServer()
	},
}

// Generate command
var generateCmd = &cobra.Command{
	Use:       "generate [certs|config|user]",
	Short:     "Generate certs or config files or a user",
	Long:      "Generate certificates, configuration files or a user required by the application.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), validGenerateArgs()),
	ValidArgs: []string{"certs", "config", "user"},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "certs":
			err := certs.CreateSelfSignedCert(
				config.Config.GetString("network.ssl-config.certificatepath"),
				config.Config.GetString("network.ssl-config.privatekeypath"))
			if err != nil {
				fmt.Println("There was an error generating the self signed certs: %w", err)
			}
		case "config":
			config.CreateDefaultConfig(config.Config)
		case "user":
			err := GenerateUser()
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("Invalid argument. Use 'certs' or 'config'.")
		}
	},
}

func validGenerateArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		validArgs := []string{"certs", "config", "user"}
		for _, valid := range validArgs {
			if args[0] == valid {
				return nil
			}
		}
		return fmt.Errorf("invalid argument: '%s'. Must be one of 'certs' or 'config' or 'user'", args[0])
	}
}

func init() {
	// Add commands to rootCmd
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(stopCmd)

	// Initialize Viper config
	cobra.OnInitialize(initConfig)
}

// generate one-of admin for login and Setup

func GenerateUser() error {
	fmt.Println("Starting to generate user")
	config.Config = config.StartViper(viper.New())

	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	username := config.Config.GetString("database.username")
	password := config.Config.GetString("database.password")
	dbAddress := config.Config.GetString("database.address")
	dbPort := config.Config.GetString("database.port")
	dbConnectionURI := fmt.Sprint("mongodb://", username, ":", password, "@", dbAddress, ":", dbPort, "/")

	fmt.Println(dbAddress)
	fmt.Println(dbConnectionURI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnectionURI))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Connected to db")

	// Ensure the client disconnects properly
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Select the collection
	coll := client.Database("packagelock").Collection("users")
	fmt.Println(coll)

	// Create a new user
	usr := structs.User{
		UserID:       uuid.New(),
		Username:     "Nick",
		Password:     "NicksPasswort",
		Groups:       []string{"Admin", "Group2"},
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		ApiKeys:      []structs.ApiKey{},
	}
	fmt.Println(usr)

	// Insert the user into the collection
	result, err := coll.InsertOne(ctx, usr)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}

// initConfig initializes Viper and configures the application
func initConfig() {
	config.Config = config.StartViper(viper.New())

	// If AppVersion is injected, set it in the configuration
	if AppVersion != "" {
		config.Config.SetDefault("general.app-version", AppVersion)
	}

	// Check and create self-signed certificates if missing
	if _, err := os.Stat(config.Config.GetString("network.ssl-config.certificatepath")); os.IsNotExist(err) {
		fmt.Println("Certificate files missing, creating new self-signed.")
		err := certs.CreateSelfSignedCert(
			config.Config.GetString("network.ssl-config.certificatepath"),
			config.Config.GetString("network.ssl-config.privatekeypath"))
		if err != nil {
			fmt.Printf("Error creating self-signed certificate: %v\n", err)
			os.Exit(1)
		}
	}
}

// startServer starts the Fiber server with appropriate configuration
func startServer() {
	pid := os.Getpid()
	err := os.WriteFile("packagelock.pid", []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		fmt.Printf("Failed to write PID file: %v\n", err)
		return
	}

	fmt.Println(config.Config.AllSettings())

	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		for {
			router := server.AddRoutes(config.Config)

			// Setup server address from config
			serverAddr := config.Config.GetString("network.fqdn") + ":" + config.Config.GetString("network.port")

			// Start server based on SSL config
			go func() {
				if config.Config.GetBool("network.ssl") {
					fmt.Printf("Starting Fiber HTTPS server at https://%s...\n", serverAddr)
					err := server.ListenAndServeTLS(
						router.Router,
						config.Config.GetString("network.ssl-config.certificatepath"),
						config.Config.GetString("network.ssl-config.privatekeypath"),
						serverAddr)
					if err != nil {
						fmt.Printf("Server error: %s\n", err)
					}
				} else {
					fmt.Printf("Starting Fiber server at %s...\n", serverAddr)
					if err := router.Router.Listen(serverAddr); err != nil {
						fmt.Printf("Server error: %s\n", err)
					}
				}
			}()

			// Handle restart or quit signals
			select {
			case <-restartChan:
				fmt.Println("Restarting Fiber server...")
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := router.Router.Shutdown(); err != nil {
					fmt.Printf("Server shutdown failed: %v\n", err)
				} else {
					fmt.Println("Server stopped.")
				}
				startServer()

			case <-quitChan:
				fmt.Println("Shutting down Fiber server...")
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := router.Router.Shutdown(); err != nil {
					fmt.Printf("Server shutdown failed: %v\n", err)
				} else {
					fmt.Println("Server stopped gracefully.")
				}
				return
			}
		}
	}()

	// Watch for config changes
	config.Config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		fmt.Println("Restarting to apply changes...")
		restartChan <- struct{}{}
	})
	config.Config.WatchConfig()

	// Block until quit signal is received
	<-quitChan
	fmt.Println("Main process exiting.")
}

func restartServer() {
	stopServer()
	fmt.Println("Restarting the Server...")
	time.Sleep(5 * time.Second)
	startServer()
}

func stopServer() {
	// Read the PID from the file using os.ReadFile
	data, err := os.ReadFile("packagelock.pid")
	if err != nil {
		fmt.Printf("Could not read PID file: %v\n", err)
		return
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Printf("Invalid PID found in file: %v\n", err)
		return
	}

	// Send SIGTERM to the process
	fmt.Printf("Stopping the server with PID: %d\n", pid)
	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		fmt.Printf("Failed to stop the server: %v\n", err)
	} else {
		fmt.Println("Server stopped.")
		// After successful stop, remove the PID file
		err = os.Remove("packagelock.pid")
		if err != nil {
			fmt.Printf("Failed to remove PID file: %v\n", err)
		} else {
			fmt.Println("PID file removed successfully.")
		}
	}
}

func main() {
	// Execute the Cobra root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
