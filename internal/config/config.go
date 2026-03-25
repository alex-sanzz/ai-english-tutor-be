package config

import "github.com/spf13/viper"

func InitConfig() (*ConfigApp, error) {
	vpr := viper.New()
	var cfg ConfigApp 
	vpr.SetConfigName("config")
	vpr.AddConfigPath("./files/")
	vpr.SetConfigType("yaml")

	// Enable reading from environment variables (or .env file)
	// With any prefix that starts with APP_
	// For example, APP_PORT will override the port value in config.yaml
	vpr.SetEnvPrefix("APP")

	// Automatically reads from environment variables
	vpr.AutomaticEnv()

	if err := vpr.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := vpr.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	

	return &cfg, nil


}