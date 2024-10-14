package cmd

import (
	"packagelock/config"
	"packagelock/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "packagelock",
	Short: "Packagelock CLI tool",
	Long:  `Packagelock CLI manages the server and other operations.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Initialize the logger
	var loggerError error
	logger.Logger, loggerError = logger.InitLogger()
	if loggerError != nil {
		// Essential APP-Part, so crash out asap
		panic(loggerError)
	}

	// Initialize the configuration
	config.InitConfig("")

	// Add subcommands to the root command
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(printRoutesCmd)
}
