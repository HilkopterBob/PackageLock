package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewSetupCmd(rootParams RootParams) *cobra.Command {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup PackageLock",
		Run: func(cmd *cobra.Command, args []string) {
			setup(rootParams.Logger)
		},
	}

	return setupCmd
}

func setup(logger *zap.Logger) {
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
