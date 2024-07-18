package cli

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

var cacheDir = ".gmail-cli"

func GetTokenDir() string {
	return path.Join(cacheDir, "tokens")
}

func initConfig() {
	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to get user home directory: %v", err))
	}

	// Create cache directory
	cacheDir = path.Join(home, cacheDir)

	// Search config in home directory with name ".gmail-cli/viper.json"
	viperConfigPath := path.Join(home, ".gmail-cli")
	viper.AddConfigPath(viperConfigPath)
	viper.SetConfigType("json")
	viper.SetConfigName("viper")

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Config file not found; using default config
			viper.SetDefault("accounts", []string{})
			// create new viper config
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		} else {
			// Config file was found but another error was produced
			cobra.CheckErr(fmt.Errorf("failed to read config file: %v", err))
		}
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
