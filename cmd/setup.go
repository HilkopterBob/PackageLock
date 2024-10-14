package cmd

import (
	"fmt"
	"os"
	"packagelock/logger"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewSetupCmd() *cobra.Command {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup PackageLock",
		Run: func(cmd *cobra.Command, args []string) {
			app := fx.New(
				logger.Module,
				fx.Invoke(runSetup),
			)

			if err := app.Start(cmd.Context()); err != nil {
				fmt.Println("Failed to start application for setup command:", err)
				os.Exit(1)
			}

			// Since runSetup runs synchronously, we can stop the app immediately
			if err := app.Stop(cmd.Context()); err != nil {
				fmt.Println("Failed to stop application after setup command:", err)
				os.Exit(1)
			}
		},
	}

	return setupCmd
}

func runSetup(logger *zap.Logger) {
	fmt.Println("Starting the PackageLock setup!")

	err := os.MkdirAll("logs/", os.ModePerm)
	if err != nil {
		logger.Fatal("Couldn't create 'logs' directory", zap.Error(err))
	}
	fmt.Println("Generated logs directory")

	err = os.MkdirAll("certs/", os.ModePerm)
	if err != nil {
		logger.Fatal("Couldn't create 'certs' directory", zap.Error(err))
	}
	fmt.Println("Generated certs directory")

	generateUnitFile(logger)
	fmt.Println("Generated systemd unit file")

	fmt.Println("Setup finished successfully!")
	os.Exit(0)
}

func generateUnitFile(logger *zap.Logger) {
	const systemdTemplate = `[Unit]
Description=PackageLock Management Server
After=network.target

[Service]
ExecStart={{.ExecStart}} start
Restart=always
User={{.User}}
Group={{.Group}}

[Install]
WantedBy=multi-user.target
`

	type unitFileData struct {
		ExecStart string
		User      string
		Group     string
	}

	execPath, err := os.Executable()
	if err != nil {
		logger.Fatal("Failed to get executable path", zap.Error(err))
	}
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		logger.Fatal("Failed to get absolute executable path", zap.Error(err))
	}

	data := unitFileData{
		ExecStart: execPath,
		User:      "your-user",  // Replace with actual user
		Group:     "your-group", // Replace with actual group
	}

	filePath := "/etc/systemd/system/packagelock.service"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Seems like you can't generate the unit file...")
		fmt.Println("Did you run this with 'sudo'? 🚀")
		logger.Fatal("Failed to create systemd unit file", zap.Error(err))
	}
	defer file.Close()

	tmpl, err := template.New("systemd").Parse(systemdTemplate)
	if err != nil {
		logger.Fatal("Failed to parse systemd template", zap.Error(err))
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		logger.Fatal("Failed to execute template", zap.Error(err))
	}

	fmt.Printf("Systemd unit file created at %s\n", filePath)
}
