package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		Drivers *Drivers  `mapstructure:"drivers"`
		Jobs    []*JobCfg `mapstructure:"jobs"`
	}
)

var (
	DefaultPath = "/etc/backup/config.yaml"
	config      Config
)

func Get() Config {
	return config
}

func Load(path string) (err error) {
	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return nil
}
