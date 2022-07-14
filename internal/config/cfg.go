package config

import (
	"backup/internal/drivers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

type (
	Config struct {
		Drivers *Drivers `mapstructure:"drivers"`
	}
	Drivers struct {
		S3  map[string]*drivers.S3DriverConfig        `mapstructure:"s3"`
		Dir map[string]*drivers.DirectoryDriverConfig `mapstructure:"dir"`
	}
)

var (
	DefaultPath string
	config      Config
)

func init() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	DefaultPath = path.Join(homePath, ".backup.yaml")
}

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
