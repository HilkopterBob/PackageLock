package cmd

import (
	"context"
	"fmt"
	"os"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/structs"
	"time"

	configPkg "packagelock/config"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/surrealdb/surrealdb.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func NewGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:       "generate [certs|config|admin]",
		Short:     "Generate certificates, configuration files, or an admin",
		Long:      "Generate certificates, configuration files, or an admin user required by the application.",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"certs", "config", "admin"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please specify an argument: certs, config, or admin.")
				return
			}
			switch args[0] {
			case "certs":
				app := fx.New(
					fx.Provide(func() string { return "Command Runner" }),
					logger.Module,
					configPkg.Module,
					certs.Module,
					fx.Invoke(runGenerateCerts),
				)

				if err := app.Start(context.Background()); err != nil {
					fmt.Println("Failed to start application for certificate generation:", err)
					os.Exit(1)
				}

				if err := app.Stop(context.Background()); err != nil {
					fmt.Println("Failed to stop application after certificate generation:", err)
					os.Exit(1)
				}
				os.Exit(0)
			case "config":
				app := fx.New(
					fx.Provide(func() string { return "Command Runner" }),
					certs.Module,
					logger.Module,
					configPkg.Module,
					fx.Invoke(runGenerateConfig),
				)

				if err := app.Start(context.Background()); err != nil {
					fmt.Println("Failed to start application for config generation:", err)
					os.Exit(1)
				}

				if err := app.Stop(context.Background()); err != nil {
					fmt.Println("Failed to stop application after config generation:", err)
					os.Exit(1)
				}
				os.Exit(0)
			case "admin":
				app := fx.New(
					fx.Provide(func() string { return "Command Runner" }),
					certs.Module,
					config.Module,
					logger.Module,
					db.Module,
					fx.Invoke(runGenerateAdmin),
				)

				if err := app.Start(context.Background()); err != nil {
					fmt.Println("Failed to start application for admin generation:", err)
					os.Exit(1)
				}

				if err := app.Stop(context.Background()); err != nil {
					fmt.Println("Failed to stop application after admin generation:", err)
					os.Exit(1)
				}
				os.Exit(0)
			default:
				fmt.Println("Invalid argument. Use 'certs', 'config', or 'admin'.")
			}
		},
	}

	return generateCmd
}

func runGenerateCerts(certGen *certs.CertGenerator, logger *zap.Logger, config *viper.Viper) {
	err := certGen.CreateSelfSignedCert(
		config.GetString("network.ssl-config.certificatepath"),
		config.GetString("network.ssl-config.privatekeypath"),
	)
	if err != nil {
		fmt.Printf("Error generating self-signed certs: %v\n", err)
		logger.Warn("Error generating self-signed certs", zap.Error(err))
	} else {
		logger.Info("Successfully generated self-signed certificates.")
		fmt.Println("Certificates generated successfully.")
	}
}

func runGenerateConfig(config *viper.Viper, logger *zap.Logger) {
	configPkg.CreateDefaultConfig(config, logger)
	logger.Info("Default configuration file created.")
	fmt.Println("Configuration file generated successfully.")
}

func runGenerateAdmin(db *db.Database, logger *zap.Logger) {
	adminPw, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		logger.Fatal("Error generating admin password", zap.Error(err))
	}

	// Hash the password for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPw), bcrypt.DefaultCost)
	if err != nil {
		logger.Fatal("Error hashing admin password", zap.Error(err))
	}

	// Admin data
	temporalAdmin := structs.User{
		UserID:       uuid.New(),
		Username:     "admin",
		Password:     string(hashedPassword),
		Groups:       []string{"Admin", "StorageAdmin", "Audit"},
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		ApiKeys:      nil,
	}

	// Insert admin
	adminInsertionData, err := db.DB.Create("user", temporalAdmin)
	if err != nil {
		logger.Fatal("Error inserting default admin into DB", zap.Error(err))
	}

	// Unmarshal data
	var createdUser structs.User
	err = surrealdb.Unmarshal(adminInsertionData, &createdUser)
	if err != nil {
		logger.Fatal("Error querying default admin", zap.Error(err))
	}

	pp.Println("Admin Username:", createdUser.Username)
	pp.Println("Admin Password:", adminPw) // Display the original password

	logger.Info("Admin user created successfully.")
}
