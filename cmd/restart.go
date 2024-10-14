package cmd

import (
	"context"
	"fmt"
	"packagelock/config"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/server"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewRestartCmd creates the restart command.
func NewRestartCmd(rootParams RootParams) *cobra.Command {
	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the running server",
		Run: func(cmd *cobra.Command, args []string) {
			// Create the Fx application
			app := fx.New(
				fx.Supply(rootParams),
				logger.Module,
				config.Module,
				db.Module,
				server.Module,
			)

			// Start the application
			if err := app.Start(context.Background()); err != nil {
				rootParams.Logger.Fatal("Failed to start application for restart", zap.Error(err))
			}

			// Perform the restart
			if err := restartApplication(app, rootParams.Logger); err != nil {
				rootParams.Logger.Fatal("Failed to restart application", zap.Error(err))
			}

			// Wait for the application to stop
			<-app.Done()

			// Stop the application
			if err := app.Stop(context.Background()); err != nil {
				rootParams.Logger.Fatal("Failed to stop application after restart", zap.Error(err))
			}
		},
	}

	return restartCmd
}

func restartApplication(app *fx.App, logger *zap.Logger) error {
	// Stop the application
	if err := app.Stop(context.Background()); err != nil {
		logger.Error("Failed to stop application", zap.Error(err))
		return err
	}

	// Wait before restarting
	time.Sleep(5 * time.Second)

	fmt.Println("Restarting the application...")
	logger.Info("Restarting the application...")

	// Start the application again
	if err := app.Start(context.Background()); err != nil {
		logger.Error("Failed to restart application", zap.Error(err))
		return err
	}

	return nil
}
