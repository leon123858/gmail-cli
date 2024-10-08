package configs

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
	tokensDir := path.Join(cacheDir, "tokens")
	// if directory not exist, create it
	if _, err := os.Stat(tokensDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tokensDir, 0700); err != nil {
			cobra.CheckErr(fmt.Errorf("failed to create token directory: %v", err))
		}
	}
	return tokensDir
}

func InitConfig() {
	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(fmt.Errorf("failed to get user home directory: %v", err))
	}

	// Create cache directory
	cacheDir = path.Join(home, cacheDir)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0700); err != nil {
			cobra.CheckErr(fmt.Errorf("failed to create cache directory: %v", err))
		}
	}

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
			viper.SetDefault("id", "")
			viper.SetDefault("secret", "")
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
