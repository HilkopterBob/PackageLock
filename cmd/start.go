package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"packagelock/config"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/server"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	restartChan = make(chan struct{})
	quitChan    = make(chan os.Signal, 1)
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		initServer(false)
	},
}

func initServer(printRoutes bool) {
	// Initialize the database
	err := db.InitDB()
	if err != nil {
		logger.Logger.Panicf("Got error from db.InitDB: %s", err)
	}

	// Start the server
	startServer(printRoutes)
}

func startServer(printRoutes bool) {
	pid := os.Getpid()
	err := os.WriteFile("packagelock.pid", []byte(fmt.Sprintf("%d", pid)), 0644)
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

			if printRoutes {
				routes := Router.Router.Stack() // Get all registered routes
				for _, route := range routes {
					for _, r := range route {
						fmt.Printf("%s %s\n", r.Method, r.Path)
					}
				}
				logger.Logger.Info("Printed all routes. Stopping the server...")
				stopServer()
				return
			}

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

				if err := Router.Router.Shutdown(); err != nil {
					logger.Logger.Warnf("Server shutdown failed: %v\n", err)
				} else {
					fmt.Println("Server stopped.")
					logger.Logger.Info("Server stopped.")
				}

				startServer(printRoutes)
				return

			case <-quitChan:
				fmt.Println("Shutting down Fiber server...")
				logger.Logger.Info("Shutting down Fiber server...")

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
		logger.Logger.Infof("Config file changed: %s", e.Name)
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
