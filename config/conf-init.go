package config

import (
	"bytes"
	"fmt"
	"io"

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
	config.SetConfigName("config")            // name of config file (without extension)
	config.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	config.AddConfigPath("/etc/packagelock/") // path to look for the config file in  etc/
	config.AddConfigPath(".")                 // optionally look for config in the working directory

	// if no config file found a default file will be Created
	// than a rescan. new_config is the same as config, but needs a different name
	// as it cont be argument return-store
	// if there is a different error -> panic & exit
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			CreateDefaultConfig(config)
			new_config := StartViper(config)
			return new_config
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	return config
}

func CreateDefaultConfig(config ConfigProvider) {
	// TODO: Add default config
	yamlExample := []byte(`
general:
  debug: true
  production: false
network:
  fqdn: 0.0.0.0
  port: 8080
  ssl: true
  ssl-config:
    redirecthttp: true
    allowselfsigned: true
    certificatepath: ./certs/testing.crt
    privatekeypath: ./certs/testing.key
  `)

	err := config.ReadConfig(bytes.NewBuffer(yamlExample))
	if err != nil {
		panic(fmt.Errorf("fatal error while reading config file: %w", err))
	}

	err_write := config.WriteConfigAs("./config.yaml")
	if err_write != nil {
		panic(fmt.Errorf("fatal error while writing config file: %w", err))
	}
}
