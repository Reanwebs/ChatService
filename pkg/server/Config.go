package server

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port       string `mapstructure:"PORT"`
	MongoURI   string `mapstructure:"MONGO_URI"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbName     string `mapstructure:"DB_NAME"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
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
