package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func NewRestartCmd() *cobra.Command {
	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the running server",
		Run: func(cmd *cobra.Command, args []string) {
			// Stop the running server
			fmt.Println("Stopping the server...")
			if err := stopServer(); err != nil {
				fmt.Println("Failed to stop the server:", err)
				os.Exit(1)
			}

			// Wait before restarting
			time.Sleep(5 * time.Second)

			// Start the server
			fmt.Println("Starting the server...")
			if err := startServer(); err != nil {
				fmt.Println("Failed to start the server:", err)
				os.Exit(1)
			}
		},
	}

	return restartCmd
}

// stopServer stops the running server by reading the PID file and sending SIGTERM.
func stopServer() error {
	// Read the PID from the file
	data, err := os.ReadFile("packagelock.pid")
	if err != nil {
		return fmt.Errorf("could not read PID file: %w", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("invalid PID found in file: %w", err)
	}

	// Send SIGTERM to the process
	fmt.Printf("Stopping the server with PID: %d\n", pid)

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find the process: %w", err)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to stop the server: %w", err)
	}

	fmt.Println("Server stopped.")

	// Remove the PID file
	err = os.Remove("packagelock.pid")
	if err != nil {
		fmt.Println("Failed to remove PID file:", err)
	} else {
		fmt.Println("PID file removed successfully.")
	}

	return nil
}

// startServer starts the server by creating a new Fx application and running it in a separate process.
func startServer() error {
	// Execute the start command as a new process
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	cmd := exec.Command(executable, "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server process: %w", err)
	}

	fmt.Printf("Server started with PID: %d\n", cmd.Process.Pid)
	return nil
}
