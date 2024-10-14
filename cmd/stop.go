package cmd

import (
	"fmt"
	"os"
	"packagelock/logger"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running server",
	Run: func(cmd *cobra.Command, args []string) {
		stopServer()
	},
}

func stopServer() {
	// Read the PID from the file
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
		logger.Logger.Warnf("Failed to stop the server: %v\n", err)
		return
	}

	fmt.Println("Server stopped.")
	logger.Logger.Info("Server stopped.")
	// Remove the PID file
	err = os.Remove("packagelock.pid")
	if err != nil {
		logger.Logger.Warnf("Failed to remove PID file: %v\n", err)
	} else {
		fmt.Println("PID file removed successfully.")
		logger.Logger.Info("PID file removed successfully.")
	}
}
