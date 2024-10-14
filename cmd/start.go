package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/db"
	"packagelock/handler"
	"packagelock/logger"
	"packagelock/server"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func NewStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
	}

	startCmd.Run = func(cmd *cobra.Command, args []string) {
		app := fx.New(
			fx.Provide(
				func() string {
					return "1.0.0" // Replace with actual version
				},
				logger.NewLogger,
				config.NewConfig,
			),
			certs.Module,
			db.Module,
			handler.Module,
			server.Module,
			fx.Invoke(func(*fiber.App) {}),
		)

		if err := app.Start(context.Background()); err != nil {
			fmt.Println("Failed to start server application:", err)
			os.Exit(1)
		}

		// Wait for interrupt signal to gracefully shutdown the server
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		if err := app.Stop(context.Background()); err != nil {
			fmt.Println("Failed to stop server application:", err)
			os.Exit(1)
		}
	}

	return startCmd
}
