package cmd

import (
	"context"
	"fmt"
	"os"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/db"
	"packagelock/handler"
	"packagelock/logger"
	"packagelock/server"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewPrintRoutesCmd() *cobra.Command {
	printRoutesCmd := &cobra.Command{
		Use:   "print-routes",
		Short: "Prints out all registered routes",
		Run: func(cmd *cobra.Command, args []string) {
			app := fx.New(
				fx.Provide(func() string { return "Command Runner" }),
				logger.Module,
				server.Module,
				handler.Module,
				config.Module,
				certs.Module,
				db.Module,
				fx.Invoke(runPrintRoutes),
			)

			if err := app.Start(context.Background()); err != nil {
				fmt.Println("Failed to start application for printing routes:", err)
				os.Exit(1)
			}

			// Since runPrintRoutes runs synchronously, we can stop the app immediately
			if err := app.Stop(context.Background()); err != nil {
				fmt.Println("Failed to stop application after printing routes:", err)
				os.Exit(1)
			}
		},
	}

	return printRoutesCmd
}

func runPrintRoutes(app *fiber.App, logger *zap.Logger) {
	routes := app.Stack() // Get all registered routes
	for _, route := range routes {
		for _, r := range route {
			fmt.Printf("%s %s\n", r.Method, r.Path)
		}
	}
	logger.Info("Printed all routes.")
	os.Exit(0)
}
