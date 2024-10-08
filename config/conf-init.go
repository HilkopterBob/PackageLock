package config

import (
	"bytes"
	"io"
	"packagelock/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config ConfigProvider

type ConfigProvider interface {
	SetConfigName(name string)
	SetConfigType(fileext string)
	AddConfigPath(path string)
	ReadInConfig() error
	OnConfigChange(run func(e fsnotify.Event))
	WatchConfig()
	WriteConfigAs(path string) error
	ReadConfig(in io.Reader) error
	AllSettings() map[string]any
	GetString(string string) string
	SetDefault(key string, value any)
	Get(key string) any
	GetBool(key string) bool
}

// TODO: How to test?
func StartViper(config ConfigProvider) ConfigProvider {
	config.SetConfigName("config") // name of config file (without extension)
	config.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	config.AddConfigPath("/app/data")
	config.AddConfigPath("/etc/packagelock/") // path to look for the config file in  etc/
	config.AddConfigPath(".")                 // optionally look for config in the working directory

	// if no config file found a default file will be Created
	// than a rescan. new_config is the same as config, but needs a different name
	// as it cont be argument return-store
	// if there is a different error -> panic & exit
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			CreateDefaultConfig(config)
			newConfig := StartViper(config)
			logger.Logger.Info("No Config found, created default Config.")
			return newConfig
		} else {
			logger.Logger.Panicf("Cannot create default config, got: %s", err)
		}
	}

	logger.Logger.Info("Successfully Created Config Manager.")
	return config
}

func CreateDefaultConfig(config ConfigProvider) {
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
    redirecthttp: true  `)

	err := config.ReadConfig(bytes.NewBuffer(yamlExample))
	if err != nil {
		logger.Logger.Panicf("Incompatible Default Config! Got: %s", err)
	}

	errWrite := config.WriteConfigAs("./config.yaml")
	if errWrite != nil {
		logger.Logger.Panicf("Cannot write config file, got: %s", errWrite)
	}
}
