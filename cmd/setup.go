package cmd

import (
	"os"
	"packagelock/logger"
	"path/filepath"
	"text/template"

	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup PackageLock",
	Run: func(cmd *cobra.Command, args []string) {
		setup()
	},
}

func setup() {
	pp.Println("Starting the PackageLock setup!")

	err := os.MkdirAll("logs/", os.ModePerm)
	if err != nil {
		pp.Printf("Couldn't create 'logs' directory. Got: %s", err)
		panic(err)
	}
	pp.Println("Generated logs directory")

	err = os.MkdirAll("certs/", os.ModePerm)
	if err != nil {
		pp.Printf("Couldn't create 'certs' directory. Got: %s", err)
		panic(err)
	}
	pp.Println("Generated certs directory")

	generateUnitFile()
	pp.Println("Generated systemd unit file")

	pp.Println("Setup finished successfully!")
}

func generateUnitFile() {
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
		logger.Logger.Panicf("Failed to get executable path: %v", err)
	}
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		logger.Logger.Panicf("Failed to get absolute executable path: %v", err)
	}

	data := unitFileData{
		ExecStart: execPath,
		User:      "your-user",  // Replace with actual user
		Group:     "your-group", // Replace with actual group
	}

	filePath := "/etc/systemd/system/packagelock.service"
	file, err := os.Create(filePath)
	if err != nil {
		pp.Println("Seems like you can't generate the unit file...")
		pp.Println("Did you run this with 'sudo'? 🚀")
		logger.Logger.Panicf("Failed to create systemd unit file: %v", err)
	}
	defer file.Close()

	tmpl, err := template.New("systemd").Parse(systemdTemplate)
	if err != nil {
		logger.Logger.Panicf("Failed to parse systemd template: %v", err)
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		logger.Logger.Panicf("Failed to execute template: %v", err)
	}

	pp.Printf("Systemd unit file created at %s\n", filePath)
}
