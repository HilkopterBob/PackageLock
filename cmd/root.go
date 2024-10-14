package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RootParams holds the dependencies for the root command.
type RootParams struct {
	fx.In

	Logger *zap.Logger
	Config *viper.Viper
}

// NewRootCmd creates the root command with injected dependencies.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "packagelock",
		Short: "Packagelock CLI tool",
		Long:  `Packagelock CLI manages the server and other operations.`,
	}

	// Add subcommands to the root command
	rootCmd.AddCommand(NewStartCmd())
	rootCmd.AddCommand(NewStopCmd())
	rootCmd.AddCommand(NewRestartCmd())
	rootCmd.AddCommand(NewSetupCmd())
	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewPrintRoutesCmd())

	return rootCmd
}

// Execute runs the root command.
func Execute(rootCmd *cobra.Command) error {
	return rootCmd.Execute()
}

// Module exports the commands as an Fx module.
var Module = fx.Options(
	fx.Provide(NewRootCmd),
	fx.Invoke(func(lc fx.Lifecycle, rootCmd *cobra.Command, logger *zap.Logger) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := rootCmd.Execute(); err != nil {
						logger.Fatal("Failed to execute root command", zap.Error(err))
					}
					os.Exit(0)
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				// Handle any cleanup if necessary
				return nil
			},
		})
	}),
)
