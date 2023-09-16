package server

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string `mapstructure:"PORT"`
	MongoURI string `mapstructure:"MONGO_URI"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Println("err reading .env viper", err)
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		log.Println("err unmarshaling config viper", err)
		return
	}

	return
}
