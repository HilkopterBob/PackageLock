package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/db"
	"packagelock/server"
	"packagelock/structs"
	"strconv"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/k0kubun/pp/v3"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/surrealdb/surrealdb.go"
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
	Use:       "generate [certs|config|admin-user]",
	Short:     "Generate certs or config files or an  admin-user",
	Long:      "Generate certificates, configuration files or an admin-user required by the application.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), validGenerateArgs()),
	ValidArgs: []string{"certs", "config", "admin"},
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
		case "admin":
			err := generateAdmin()
			if err != nil {
				// FIXME: Error Handling
				// FIXME: Logging! Because: Invocation of admin creation should be logged!
				panic(err)
			}
		default:
			fmt.Println("Invalid argument. Use 'certs' or 'config' or 'admin'.")
		}
	},
}

func validGenerateArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		validArgs := []string{"certs", "config", "admin"}
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

	initConfig()
	err := db.InitDB()
	if err != nil {
		// FIXME: error Handling
		// FIXME: LOGGING!
		panic(err)
	}
}

// generate admin for login and Setup
func generateAdmin() error {
	adminPw, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		// FIXME: error Handling
		panic(err)
	}

	// Admin Data
	TemporalAdmin := structs.User{
		UserID:       uuid.New(),
		Username:     "admin",
		Password:     adminPw,
		Groups:       []string{"Admin", "StorageAdmin", "Audit"},
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		ApiKeys:      nil,
	}

	// Insert Admin
	adminInsertionData, err := db.DB.Create("user", TemporalAdmin)
	if err != nil {
		panic(err)
	}

	// Unmarshal data
	var createdUser structs.User
	err = surrealdb.Unmarshal(adminInsertionData, &createdUser)
	if err != nil {
		pp.Println(adminInsertionData)
		panic(err)
	}

	pp.Println(createdUser.Username)
	pp.Println(createdUser.Password)
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
			Router := server.AddRoutes(config.Config)

			// Setup server address from config
			serverAddr := config.Config.GetString("network.fqdn") + ":" + config.Config.GetString("network.port")

			// Start server based on SSL config
			go func() {
				if config.Config.GetBool("network.ssl") {
					fmt.Printf("Starting Fiber HTTPS server at https://%s...\n", serverAddr)
					err := server.ListenAndServeTLS(
						Router.Router,
						config.Config.GetString("network.ssl-config.certificatepath"),
						config.Config.GetString("network.ssl-config.privatekeypath"),
						serverAddr)
					if err != nil {
						fmt.Printf("Server error: %s\n", err)
					}
				} else {
					fmt.Printf("Starting Fiber server at %s...\n", serverAddr)
					if err := Router.Router.Listen(serverAddr); err != nil {
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
				if err := Router.Router.Shutdown(); err != nil {
					fmt.Printf("Server shutdown failed: %v\n", err)
				} else {
					fmt.Println("Server stopped.")
				}
				startServer()

			case <-quitChan:
				fmt.Println("Shutting down Fiber server...")
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := Router.Router.Shutdown(); err != nil {
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
