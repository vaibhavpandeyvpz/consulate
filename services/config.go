package services

import (
	"flag"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"sync"
)

type Config struct {
	Database struct {
		Path string `mapstructure:"path"`
	} `config:"database"`
	Exotel struct {
		Domain            string `mapstructure:"domain"`
		AccountSid        string `mapstructure:"account_sid"`
		ApiKey            string `mapstructure:"api_key"`
		ApiToken          string `mapstructure:"api_token"`
		CallerId          string `mapstructure:"caller_id"`
		StatusCallbackUrl string `mapstructure:"status_callback_url"`
	} `config:"exotel"`
	Recipients []string `config:"recipients"`
	Slack      struct {
		BotToken string `mapstructure:"bot_token"`
		Channels struct {
			Enquiries string `mapstructure:"enquiries"`
		} `mapstructure:"channels"`
		SigningSecret string `mapstructure:"signing_secret"`
	} `config:"slack"`
}

var _config Config
var _configCreator sync.Once
var _configPath = flag.String("config", "config.yml", "path to the config file")

func GetConfig() Config {
	_configCreator.Do(func() {
		config.WithOptions(config.ParseEnv)
		config.AddDriver(yaml.Driver)

		err := config.LoadFiles(*_configPath)
		if err != nil {
			panic(err)
		}

		err = config.Decode(&_config)
		if err != nil {
			panic(err)
		}
	})

	return _config
}
