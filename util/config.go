package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	UploadPath    string `mapstructure:"RAND_PATH"`
	ServerAddress string `mapstructure:"RAND_LISTEN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath("$HOME/")
	viper.SetConfigName(".randrc")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
