package main

import (
	"context"
	"fmt"
	"os"
	"packagelock/certs"
	"packagelock/cmd"
	"packagelock/config"
	"packagelock/db"
	"packagelock/handler"
	"packagelock/logger"
	"packagelock/server"

	"go.uber.org/fx"
)

var AppVersion string // Version injected with ldflags

func main() {
	app := fx.New(
		fx.Provide(
			func() string {
				return AppVersion
			},
			logger.NewLogger,
			config.NewConfig,
		),
		certs.Module,
		db.Module,      // Include the database module
		handler.Module, // Include the handlers module
		server.Module,  // Include the server module
		cmd.Module,     // Include the commands module
	)

	if err := app.Start(context.Background()); err != nil {
		fmt.Println("Failed to start application:", err)
		os.Exit(1)
	}

	// Wait for the application to be signaled to exit
	<-app.Done()

	if err := app.Stop(context.Background()); err != nil {
		fmt.Println("Failed to stop application:", err)
		os.Exit(1)
	}
}
