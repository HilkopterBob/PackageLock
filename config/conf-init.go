package config

import (
	"bytes"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func StartViper() {
	viper.SetConfigName("config")            // name of config file (without extension)
	viper.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/packagelock/") // path to look for the config file in  etc/
	viper.AddConfigPath(".")                 // optionally look for config in the working directory

	// if no config file found a default file will be Created
	// than a rescan
	// if there is a different error -> panic & exit
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			CreateDefaultConfig()
			StartViper()
			return
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
}

func CreateDefaultConfig() {
	// TODO: Add default config
	yamlExample := []byte(`
general:
  debug: True
  production: False
  Port: 8080

Network:
  FQDN: "packagelock.company.com"
  ForceHTTP: False
  SSL:
    CertificatePath: "/etc/packagelock/ssl/cert.pem"
    PrivateKeyPath: "/etc/packagelock/ssl/privkey.pem"
    AllowSelfSigned: False
  `)

	err := viper.ReadConfig(bytes.NewBuffer(yamlExample))
	if err != nil {
		panic(fmt.Errorf("fatal error while reading config file: %w", err))
	}

	err_write := viper.WriteConfigAs("./config.yaml")
	if err_write != nil {
		panic(fmt.Errorf("fatal error while writing config file: %w", err))
	}
}
