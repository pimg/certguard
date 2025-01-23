package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	configFileName = "config"
	envPrefix      = "CG"
)

type ViperConfig struct {
	*viper.Viper
	cfg *Config
}

func NewViperConfig() *ViperConfig {
	return &ViperConfig{viper.New(), New()}
}

func (v *ViperConfig) Config() *Config {
	return v.cfg
}

func (v *ViperConfig) InitializeConfig() error {
	// Set the base name of the config file, without the file extension.
	v.SetConfigName(configFileName)
	v.SetConfigType("yaml")

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.config/certguard")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the setFlags function
	v.AutomaticEnv()

	v.cfg.Theme.Name = v.GetString("config.theme.name")
	v.cfg.Log.Debug = v.GetBool("config.log.debug")
	v.cfg.Log.Directory = v.GetString("config.log.file")

	return nil
}
