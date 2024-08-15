package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"packagelock/config"
	"packagelock/server"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Data structs

// TODO: support for multiple network adapters.

func main() {
	Config := config.StartViper(viper.New())
	fmt.Println(Config.AllSettings())

	// Channel to signal the restart
	restartChan := make(chan struct{})
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		for {
			router := server.AddRoutes()
			serverAddr := Config.GetString("network.fqdn") + ":" + Config.GetString("network.port")
			srv := &http.Server{
				Addr:    serverAddr,
				Handler: router.Router.Handler(),
			}

			go func() {
				fmt.Printf("Starting server at %s...\n", serverAddr)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					fmt.Printf("Server error: %s\n", err)
				}
			}()

			// Wait for either a restart signal or termination signal
			select {
			case <-restartChan:
				fmt.Println("Restarting server...")

				// Gracefully shutdown the server
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					fmt.Printf("Server shutdown failed: %v\n", err)
				} else {
					fmt.Println("Server stopped.")
				}

			case <-quitChan:
				fmt.Println("Shutting down server...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					fmt.Printf("Server shutdown failed: %v\n", err)
				} else {
					fmt.Println("Server stopped gracefully.")
				}
				return
			}
		}
	}()

	// Watch for configuration changes
	Config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		fmt.Println("Restarting to apply changes...")
		restartChan <- struct{}{} // Send signal to restart the server
	})
	Config.WatchConfig()

	// Block until quit signal is received
	<-quitChan
	fmt.Println("Main process exiting.")
}
