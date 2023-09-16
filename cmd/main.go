package main

import (
	"chat/pkg/server"
	"log"
)

func main() {
	config, err := server.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ginServer, err := server.InitializeApi(config)
	if err != nil {
		log.Fatalln(err)
	}

	server.Start(ginServer, config)
}
