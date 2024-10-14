package cmd

import (
	"context"
	"fmt"
	"packagelock/server"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewPrintRoutesCmd(rootParams RootParams) *cobra.Command {
	printRoutesCmd := &cobra.Command{
		Use:   "print-routes",
		Short: "Prints out all registered routes",
		Run: func(cmd *cobra.Command, args []string) {
			app := fx.New(
				fx.Supply(rootParams),
				server.Module,
				fx.Invoke(printRoutes),
			)

			if err := app.Start(context.Background()); err != nil {
				rootParams.Logger.Fatal("Failed to start application for printing routes", zap.Error(err))
			}

			if err := app.Stop(context.Background()); err != nil {
				rootParams.Logger.Fatal("Failed to stop application after printing routes", zap.Error(err))
			}
		},
	}

	return printRoutesCmd
}

func printRoutes(server *fiber.App, logger *zap.Logger) {
	routes := server.Stack() // Get all registered routes
	for _, route := range routes {
		for _, r := range route {
			fmt.Printf("%s %s\n", r.Method, r.Path)
		}
	}
	logger.Info("Printed all routes.")
}
