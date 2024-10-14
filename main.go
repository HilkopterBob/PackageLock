package main

import (
	"fmt"
	"os"
	"packagelock/cmd"
	"packagelock/config"
)

var AppVersion string // Version injected with ldflags

func main() {
	// Initialize the configuration
	config.InitConfig(AppVersion)

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
