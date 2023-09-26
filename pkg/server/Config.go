package server

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	AuthClientURL string `mapstructure:"Auth_SRV"`
	MongoURI      string `mapstructure:"MONGO_URI"`
	DbHost        string `mapstructure:"DB_HOST"`
	DbName        string `mapstructure:"DB_NAME"`
	DbUser        string `mapstructure:"DB_USER"`
	DbPort        string `mapstructure:"DB_PORT"`
	DbPassword    string `mapstructure:"DB_PASSWORD"`
	JwtKey        string `mapstructure:"JWT_SECRET_KEY"`
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
