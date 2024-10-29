package config

import (
	"bytes"
	"context"
	"os"
	"packagelock/certs"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/codes" // Import for setting span status
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ConfigParams struct {
	fx.In

	Lifecycle     fx.Lifecycle
	Logger        *zap.Logger
	AppVersion    string
	CertGenerator *certs.CertGenerator
	Tracer        trace.Tracer // Injected Tracer from OpenTelemetry
}

func NewConfig(params ConfigParams) (*viper.Viper, error) {
	// Start a new span for the configuration initialization
	_, span := params.Tracer.Start(context.Background(), "Configuration Initialization")
	defer span.End()

	config := viper.New()
	config.SetDefault("general.app-version", params.AppVersion)
	config.SetConfigName("config") // Name of config file (without extension)
	config.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	config.AddConfigPath("/app/data")
	config.AddConfigPath("/etc/packagelock/") // Path to look for the config file in etc/
	config.AddConfigPath(".")                 // Optionally look for config in the working directory

	// Add attributes to the span
	span.SetAttributes()

	// Read the configuration file
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create a default config if none is found
			params.Logger.Info("No config file found. Creating default configuration.")
			span.AddEvent("Config file not found. Creating default configuration.")

			CreateDefaultConfig(config, params.Logger)

			// Attempt to read the config again
			if err := config.ReadInConfig(); err != nil {
				params.Logger.Panic("Cannot read config after creating default config", zap.Error(err))
				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to read config after creating default")
				return nil, err
			}
			params.Logger.Info("Default config created and loaded successfully.")
			span.AddEvent("Default config created and loaded successfully.")
		} else {
			params.Logger.Panic("Cannot read config", zap.Error(err))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to read config")
			return nil, err
		}
	}

	params.Logger.Info("Successfully created Config Manager.")
	span.AddEvent("Config Manager initialized successfully.")

	// Check and create self-signed certificates if missing
	certPath := config.GetString("network.ssl-config.certificatepath")
	keyPath := config.GetString("network.ssl-config.privatekeypath")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		params.Logger.Info("Certificate files missing. Creating new self-signed certificates.")
		span.AddEvent("Certificate files missing. Creating new self-signed certificates.")

		err := params.CertGenerator.CreateSelfSignedCert(certPath, keyPath)
		if err != nil {
			params.Logger.Panic("Error creating self-signed certificate", zap.Error(err))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to create self-signed certificate")
			return nil, err
		}

		params.Logger.Info("Self-signed certificates created successfully.")
		span.AddEvent("Self-signed certificates created successfully.")
	}

	// Set up configuration change watching
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start a new span for setting up config change watcher
			_, watchSpan := params.Tracer.Start(ctx, "Config Watcher Setup")
			defer watchSpan.End()

			// Watch for configuration changes
			config.OnConfigChange(func(e fsnotify.Event) {
				params.Logger.Info("Config file changed", zap.String("file", e.Name))
				watchSpan.AddEvent("Configuration file changed")
				// Handle configuration change if necessary
			})
			config.WatchConfig()
			params.Logger.Info("Started watching configuration changes.")
			watchSpan.AddEvent("Started watching configuration changes.")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Handle any cleanup if necessary
			params.Logger.Info("Stopping configuration watcher.")
			span.AddEvent("Configuration watcher stopped.")
			return nil
		},
	})

	params.Logger.Info("Configuration initialized successfully.")
	span.SetStatus(codes.Ok, "Configuration initialized successfully.")
	span.AddEvent("Configuration initialized successfully.")

	return config, nil
}

// CreateDefaultConfig generates a default configuration file.
func CreateDefaultConfig(config *viper.Viper, logger *zap.Logger) {
	yamlExample := []byte(`
general:
  debug: true
  production: false
	monitoring: true
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
