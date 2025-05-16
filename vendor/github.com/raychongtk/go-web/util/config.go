package util

import (
	"github.com/spf13/viper"
	"log"
)

func Load(name string, path string, config interface{}) {
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	marshalErr := viper.Unmarshal(&config)
	if marshalErr != nil {
		log.Fatalln(marshalErr)
	}
}
