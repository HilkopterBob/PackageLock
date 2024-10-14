package cmd

import (
	"fmt"
	"packagelock/logger"
	"time"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the running server",
	Run: func(cmd *cobra.Command, args []string) {
		restartServer()
	},
}

func restartServer() {
	stopServer()
	time.Sleep(5 * time.Second)
	fmt.Println("Restarting the server...")
	logger.Logger.Info("Restarting the server...")
	startServer(false)
}
