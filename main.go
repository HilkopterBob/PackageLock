package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/server"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// Generate command
var generateCmd = &cobra.Command{
	Use:       "generate [certs|config]",
	Short:     "Generate certs or config files",
	Long:      "Generate certificates or configuration files required by the application.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), validGenerateArgs()), // Expect exactly one argument: either "certs" or "config"
	ValidArgs: []string{"certs", "config"},                             // Restrict arguments to these options
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
		default:
			fmt.Println("Invalid argument. Use 'certs' or 'config'.")
		}
	},
}

func validGenerateArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		validArgs := []string{"certs", "config"}
		for _, valid := range validArgs {
			if args[0] == valid {
				return nil
			}
		}
		return fmt.Errorf("invalid argument: '%s'. Must be one of 'certs' or 'config'", args[0])
	}
}

func init() {
	// Add commands to rootCmd
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(generateCmd)

	// Initialize Viper config
	cobra.OnInitialize(initConfig)
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
					if err := server.ListenAndServeTLS(
						router.Router,
						config.Config.GetString("network.ssl-config.certificatepath"),
						config.Config.GetString("network.ssl-config.privatekeypath"),
						serverAddr); err != nil {
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

func main() {
	// Execute the Cobra root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
