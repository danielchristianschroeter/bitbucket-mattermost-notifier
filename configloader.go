package main

import (
	"log"

	"github.com/spf13/viper"
)

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() {
	viper.SetDefault("LISTEN_PORT", "1337")
	viper.SetDefault("MATTERMOST_USERNAME", "Bitbucket")
	viper.SetDefault("MESSAGE_ICON_EMOJI", ":exclamation:")
	viper.SetDefault("INFO_COLOR", "3498db")
	viper.SetDefault("DECLINED_COLOR", "e74c3c")
	viper.SetDefault("WARNING_COLOR", "f1c40f")
	viper.SetDefault("SUCCESS_COLOR", "2ecc71")
	viper.SetConfigName("config")                              // name of config file (without extension)
	viper.SetConfigType("env")                                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/bitbucket-mattermost-notifier/") // path to look for the config file in
	viper.AddConfigPath(".")                                   // optionally look for config in the working directory
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatal("Fatal error config file: %w", err)
		}
	}
	if !checkParams("MATTERMOST_WEBHOOKURL", "MATTERMOST_CHANNEL") {
		log.Fatal("MATTERMOST_WEBHOOKURL or MATTERMOST_CHANNEL can NOT be empty")
	}
}

func checkParams(params ...string) bool {
	for _, param := range params {
		if !viper.IsSet(param) {
			return false
		}
	}
	return true
}

func GetConfigValue(key string) string {
	LoadConfig()
	return viper.GetString(key)
}

func GetConfigBoolValue(key string) bool {
	LoadConfig()
	return viper.GetBool(key)
}
