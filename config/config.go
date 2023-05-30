package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type TwitchConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type ServerConfig struct {
	Port    string
	BaseURL string
}

type RouterConfig struct {
	AuthURL   string
	NoAuthURL string
}

type Config struct {
	ServerConfig
	TwitchConfig
	RouterConfig
}

func NewConfig() (Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read configuration file: %v", err)
	}

	serverConfig := ServerConfig{
		Port:    viper.GetString("serverConfig.port"),
		BaseURL: viper.GetString("serverConfig.baseURL"),
	}
	routerConfig := RouterConfig{
		AuthURL:   viper.GetString("routerConfig.authURL"),
		NoAuthURL: viper.GetString("routerConfig.noAuthURL"),
	}
	twitchConfig := TwitchConfig{
		ClientID:     viper.GetString("twitchConfig.clientID"),
		ClientSecret: viper.GetString("twitchConfig.clientSecret"),
		RedirectURL:  viper.GetString("twitchConfig.redirectURL"),
		Scopes:       viper.GetStringSlice("twitchConfig.scopes"),
	}

	config = Config{
		ServerConfig: serverConfig,
		RouterConfig: routerConfig,
		TwitchConfig: twitchConfig,
	}

	fmt.Printf("config.Port: %s \n", config.Port)

	return config, nil
}
