package cmd

import (
	"fmt"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/structs"
	"time"

	"github.com/google/uuid"
	"github.com/k0kubun/pp/v3"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/surrealdb/surrealdb.go"
)

var generateCmd = &cobra.Command{
	Use:       "generate [certs|config|admin]",
	Short:     "Generate certificates, configuration files, or an admin",
	Long:      "Generate certificates, configuration files, or an admin user required by the application.",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"certs", "config", "admin"},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "certs":
			err := certs.CreateSelfSignedCert(
				config.Config.GetString("network.ssl-config.certificatepath"),
				config.Config.GetString("network.ssl-config.privatekeypath"))
			if err != nil {
				fmt.Printf("There was an error generating the self-signed certs: %v\n", err)
				logger.Logger.Warnf("There was an error generating the self-signed certs: %v", err)
			}
		case "config":
			config.CreateDefaultConfig(config.Config)
		case "admin":
			err := generateAdmin()
			if err != nil {
				logger.Logger.Panicf("Failed to generate default admin, got: %v", err)
			}
		default:
			fmt.Println("Invalid argument. Use 'certs', 'config', or 'admin'.")
		}
	},
}

func generateAdmin() error {
	// Initialize the database
	err := db.InitDB()
	if err != nil {
		logger.Logger.Panicf("Got error from db.InitDB: %v", err)
	}

	adminPw, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		logger.Logger.Panicf("Got error while generating admin password: %v", err)
	}

	// Admin data
	temporalAdmin := structs.User{
		UserID:       uuid.New(),
		Username:     "admin",
		Password:     adminPw,
		Groups:       []string{"Admin", "StorageAdmin", "Audit"},
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		ApiKeys:      nil,
	}

	// Insert admin
	adminInsertionData, err := db.DB.Create("user", temporalAdmin)
	if err != nil {
		logger.Logger.Panicf("Got error while inserting default admin into DB: %v", err)
	}

	// Unmarshal data
	var createdUser structs.User
	err = surrealdb.Unmarshal(adminInsertionData, &createdUser)
	if err != nil {
		logger.Logger.Panicf("Got error while querying default admin: %v", err)
	}

	pp.Println("Admin Username:", createdUser.Username)
	pp.Println("Admin Password:", createdUser.Password)
	return nil
}
