package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/server"
	"packagelock/structs"
	"path/filepath"
	"strconv"
	"syscall"
	"text/template"
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

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup PackageLock",
	Run: func(cmd *cobra.Command, args []string) {
		setup()
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
				logger.Logger.Warnf("There was an error generating the self signed certs: %s", err)
			}
		case "config":
			config.CreateDefaultConfig(config.Config)
		case "admin":
			err := generateAdmin()
			if err != nil {
				logger.Logger.Panicf("Failed to generate default admin, got: %s", err)
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

func setup() {
	pp.Println("Starting The PackageLock Setup!")

	err := os.MkdirAll("logs/", os.ModePerm)
	if err != nil {
		pp.Printf("Couldn't create 'logs' directory. Got: %s", err)
		panic(err)
	}
	pp.Println("Generated Logs directory")

	err = os.MkdirAll("certs/", os.ModePerm)
	if err != nil {
		pp.Printf("Couldn't create 'logs' directory. Got: %s", err)
		panic(err)
	}
	pp.Println("Generated certs directory")

	generateUnitFile()
	pp.Println("Generated Unit File")

	pp.Println("Setup finished successfully!")
}

func generateUnitFile() {
	// SystemdTemplate defines the systemd unit file structure
	const SystemdTemplate = `[Unit]
Description=PackageLock Management Server
After=network.target

[Service]
ExecStart={{.ExecStart}} start
Restart=always
User={{.User}}
Group={{.Group}}

[Install]
WantedBy=multi-user.target
`

	// UnitFileData holds the dynamic data for the systemd unit file
	type UnitFileData struct {
		ExecStart string
		User      string
		Group     string
	}

	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		logger.Logger.Panicf("failed to get executable path: %w", err)
	}
	// Convert to an absolute path
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		logger.Logger.Panicf("failed to get absolute executable path: %w", err)
	}

	// Define the data to be injected into the unit file
	data := UnitFileData{
		ExecStart: execPath,     // The path of the Go binary
		User:      "your-user",  // Replace with your actual user
		Group:     "your-group", // Replace with your actual group
	}

	// Open the systemd unit file for writing (requires sudo permission)
	filePath := "/etc/systemd/system/packagelock.service"
	file, err := os.Create(filePath)
	if err != nil {
		pp.Println("Seems like you cant generate the Unit File...")
		pp.Println("Did you ran this with 'sudo'?🚀")
		logger.Logger.Panicf("failed to create systemd unit file: %w", err)
	}
	defer file.Close()

	// Parse and execute the systemd template
	tmpl, err := template.New("systemd").Parse(SystemdTemplate)
	if err != nil {
		logger.Logger.Panicf("failed to parse systemd template: %w", err)
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		logger.Logger.Panicf("failed to execute template: %w", err)
	}

	pp.Printf("Systemd unit file created at %s\n", filePath)
}

// INFO: init is ran everytime a cobra comand gets used.
// It does not init the Server!
// It only inits cobra!
func init() {
	// Add commands to rootCmd
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(setupCmd)

	// Declare the Logger into global logger.Logger
	// Init here so commands can be logget to!
	var loggerError error
	logger.Logger, loggerError = logger.InitLogger()
	if loggerError != nil {
		// INFO: Essential APP-Part, so crash out asap
		panic(loggerError)
	}
}

// generate admin for login and Setup
func generateAdmin() error {
	adminPw, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		logger.Logger.Panicf("Got error while generating ADmin Password: %s", err)
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
		logger.Logger.Panicf("Got error while inserting Default Admin into DB: %s", err)
	}

	// Unmarshal data
	var createdUser structs.User
	err = surrealdb.Unmarshal(adminInsertionData, &createdUser)
	if err != nil {
		logger.Logger.Panicf("Got error while querring Default Admin: %s", err)
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
			logger.Logger.Panicf("Error creating self-signed certificate: %v\n", err)
		}
	}
}

// Initializes everything that is needed for the Server
// to run
func initServer() {
	initConfig()
	err := db.InitDB()
	if err != nil {
		logger.Logger.Panicf("Got error from db.InitDB: %s", err)
	}

	// after init run Server
	startServer()
}

// startServer starts the Fiber server with appropriate configuration
func startServer() {
	pid := os.Getpid()
	err := os.WriteFile("packagelock.pid", []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		logger.Logger.Panicf("Failed to write PID file: %v\n", err)
		return
	}

	if config.Config.GetString("general.production") == "false" {
		logger.Logger.Debug(config.Config.AllSettings())
	}

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

					logger.Logger.Infof("Starting Fiber HTTPS server at https://%s...\n", serverAddr)

					err := server.ListenAndServeTLS(
						Router.Router,
						config.Config.GetString("network.ssl-config.certificatepath"),
						config.Config.GetString("network.ssl-config.privatekeypath"),
						serverAddr)
					if err != nil {
						logger.Logger.Panicf("Server error: %s\n", err)
					}
				} else {
					logger.Logger.Infof("Starting Fiber server at %s...\n", serverAddr)

					if err := Router.Router.Listen(serverAddr); err != nil {
						logger.Logger.Panicf("Server error: %s\n", err)
					}
				}
			}()

			// Handle restart or quit signals
			select {
			case <-restartChan:

				fmt.Println("Restarting Fiber server...")
				logger.Logger.Info("Restarting Fiber server...")

				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := Router.Router.Shutdown(); err != nil {
					logger.Logger.Warnf("Server shutdown failed: %v\n", err)
				} else {
					// TODO: add Reason for restart/Stoping
					fmt.Println("Server stopped.")
					logger.Logger.Info("Server stopped.")
				}

				startServer()

			case <-quitChan:

				// TODO: add Reason fro Stopping
				fmt.Println("Shutting down Fiber server...")
				logger.Logger.Info("Shutting down Fiber server...")
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := Router.Router.Shutdown(); err != nil {
					logger.Logger.Warnf("Server shutdown failed: %v\n", err)
				} else {
					fmt.Println("Server stopped gracefully.")
					logger.Logger.Info("Server stopped gracefully.")
				}
				return
			}
		}
	}()

	// Watch for config changes
	config.Config.OnConfigChange(func(e fsnotify.Event) {
		logger.Logger.Infof("Config file changed:", e.Name)
		logger.Logger.Info("Restarting to apply changes...")
		fmt.Println("Restarting to apply changes...")
		restartChan <- struct{}{}
	})
	config.Config.WatchConfig()

	// Block until quit signal is received
	<-quitChan
	logger.Logger.Info("Main process exiting.")
	fmt.Println("Main process exiting.")
}

func restartServer() {
	stopServer()
	fmt.Println("Restarting the Server...")
	logger.Logger.Info("Restarting the Server...")
	time.Sleep(5 * time.Second)
	startServer()
}

func stopServer() {
	// Read the PID from the file using os.ReadFile
	data, err := os.ReadFile("packagelock.pid")
	if err != nil {
		logger.Logger.Panicf("Could not read PID file: %v\n", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		logger.Logger.Panicf("Invalid PID found in file: %v\n", err)
	}

	// Send SIGTERM to the process
	fmt.Printf("Stopping the server with PID: %d\n", pid)
	logger.Logger.Infof("Stopping the server with PID: %d\n", pid)
	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		logger.Logger.Warn("Failed to stop the server: %v\n", err)
		return
	}

	fmt.Println("Server stopped.")
	logger.Logger.Info("Server stopped.")
	// After successful stop, remove the PID file
	err = os.Remove("packagelock.pid")
	if err != nil {
		logger.Logger.Warnf("Failed to remove PID file: %v\n", err)
	} else {
		fmt.Println("PID file removed successfully.")
		logger.Logger.Info("PID file removed successfully.")
	}
}

func main() {
	// Execute the Cobra root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
