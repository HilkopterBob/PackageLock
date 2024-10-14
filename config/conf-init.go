package config

import (
	"bytes"
	"context"
	"os"
	"packagelock/certs"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ConfigParams struct {
	fx.In

	Lifecycle     fx.Lifecycle
	Logger        *zap.Logger
	AppVersion    string
	CertGenerator *certs.CertGenerator
}

func NewConfig(params ConfigParams) (*viper.Viper, error) {
	config := viper.New()
	config.SetDefault("general.app-version", params.AppVersion)
	config.SetConfigName("config") // Name of config file (without extension)
	config.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	config.AddConfigPath("/app/data")
	config.AddConfigPath("/etc/packagelock/") // Path to look for the config file in etc/
	config.AddConfigPath(".")                 // Optionally look for config in the working directory

	// Read the configuration file
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create a default config if none is found
			CreateDefaultConfig(config, params.Logger)
			// Attempt to read the config again
			if err := config.ReadInConfig(); err != nil {
				params.Logger.Panic("Cannot read config after creating default config", zap.Error(err))
				return nil, err
			}
			params.Logger.Info("No config found, created default config.")
		} else {
			params.Logger.Panic("Cannot read config", zap.Error(err))
			return nil, err
		}
	}

	params.Logger.Info("Successfully created Config Manager.")

	// Check and create self-signed certificates if missing
	certPath := config.GetString("network.ssl-config.certificatepath")
	keyPath := config.GetString("network.ssl-config.privatekeypath")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		params.Logger.Info("Certificate files missing, creating new self-signed.")
		err := params.CertGenerator.CreateSelfSignedCert(certPath, keyPath)
		if err != nil {
			params.Logger.Panic("Error creating self-signed certificate", zap.Error(err))
			return nil, err
		}
	}

	// Set up configuration change watching
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// Watch for configuration changes
			config.OnConfigChange(func(e fsnotify.Event) {
				params.Logger.Info("Config file changed", zap.String("file", e.Name))
				// Handle configuration change if necessary
			})
			config.WatchConfig()
			return nil
		},
		OnStop: func(context.Context) error {
			// Handle any cleanup if necessary
			return nil
		},
	})

	return config, nil
}

// CreateDefaultConfig generates a default configuration file.
func CreateDefaultConfig(config *viper.Viper, logger *zap.Logger) {
	yamlExample := []byte(`
general:
  debug: true
  production: false
database:
  address: 127.0.0.1
  port: 8000
  username: root
  password: root
network:
  fqdn: 0.0.0.0
  port: 8080
  ssl: true
  ssl-config:
    allowselfsigned: true
    certificatepath: ./certs/testing.crt
    privatekeypath: ./certs/testing.key
    redirecthttp: true
`)

	// Read the default configuration from the YAML example
	err := config.ReadConfig(bytes.NewBuffer(yamlExample))
	if err != nil {
		logger.Panic("Incompatible default config", zap.Error(err))
	}

	// Write the default configuration to a file
	errWrite := config.WriteConfigAs("./config.yaml")
	if errWrite != nil {
		logger.Panic("Cannot write config file", zap.Error(errWrite))
	}
}

// Module exports the config module.
var Module = fx.Options(
	fx.Provide(NewConfig),
)
