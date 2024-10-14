package cmd

import (
	"context"
	"fmt"
	"packagelock/certs"
	"packagelock/config"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/structs"
	"time"

	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/surrealdb/surrealdb.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func NewGenerateCmd(rootParams RootParams) *cobra.Command {
	generateCmd := &cobra.Command{
		Use:       "generate [certs|config|admin]",
		Short:     "Generate certificates, configuration files, or an admin",
		Long:      "Generate certificates, configuration files, or an admin user required by the application.",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"certs", "config", "admin"},
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "certs":
				app := fx.New(
					fx.Supply(rootParams),
					logger.Module,
					config.Module,
					certs.Module,
					fx.Invoke(func(certGen *certs.CertGenerator, logger *zap.Logger, config *viper.Viper) {
						err := certGen.CreateSelfSignedCert(
							config.GetString("network.ssl-config.certificatepath"),
							config.GetString("network.ssl-config.privatekeypath"),
						)
						if err != nil {
							fmt.Printf("Error generating self-signed certs: %v\n", err)
							logger.Warn("Error generating self-signed certs", zap.Error(err))
						}
					}),
				)

				if err := app.Start(context.Background()); err != nil {
					rootParams.Logger.Fatal("Failed to start application for certificate generation", zap.Error(err))
				}

				if err := app.Stop(context.Background()); err != nil {
					rootParams.Logger.Fatal("Failed to stop application after certificate generation", zap.Error(err))
				}
			case "config":
				config.CreateDefaultConfig(rootParams.Config, rootParams.Logger)
			case "admin":
				app := fx.New(
					fx.Supply(rootParams),
					db.Module,
					fx.Invoke(generateAdmin),
				)

				if err := app.Start(context.Background()); err != nil {
					rootParams.Logger.Fatal("Failed to start application for admin generation", zap.Error(err))
				}

				if err := app.Stop(context.Background()); err != nil {
					rootParams.Logger.Fatal("Failed to stop application after admin generation", zap.Error(err))
				}
			default:
				fmt.Println("Invalid argument. Use 'certs', 'config', or 'admin'.")
			}
		},
	}

	return generateCmd
}

// Add the missing generateAdmin function
func generateAdmin(db *db.Database, logger *zap.Logger) {
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

	fmt.Println("Admin Username:", createdUser.Username)
	fmt.Println("Admin Password:", adminPw) // Display the original password

	// For security, consider providing the password securely rather than printing
}
