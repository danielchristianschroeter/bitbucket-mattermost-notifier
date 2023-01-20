package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() {
	viper.SetDefault("MESSAGE_ICON_EMOJI", ":exclamation:")
	viper.SetDefault("INFO_COLOR", "3498db")
	viper.SetDefault("DECLINED_COLOR", "e74c3c")
	viper.SetDefault("WARNING_COLOR", "f1c40f")
	viper.SetDefault("SUCCESS_COLOR", "2ecc71")
	viper.SetEnvPrefix("BMN")                                  // prefix while reading from the environment variables
	viper.SetConfigName("config")                              // name of config file (without extension)
	viper.SetConfigType("env")                                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/bitbucket-mattermost-notifier/") // path to look for the config file in
	viper.AddConfigPath(".")                                   // optionally look for config in the working directory
	viper.AutomaticEnv()
	if !checkParams("LISTEN_PORT", "MATTERMOST_WEBHOOKURL", "MATTERMOST_USERNAME", "MATTERMOST_CHANNEL") {
		log.Fatal("LISTEN_PORT, MATTERMOST_WEBHOOKURL, MATTERMOST_USERNAME or MATTERMOST_CHANNEL can NOT be empty")
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
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
