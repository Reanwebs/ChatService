package main

import (
	"chat/pkg/client"
	"chat/pkg/server"
	"log"
)

func main() {
	config, err := server.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	grpcClient, err := client.InitializeClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	ginServer, err := server.InitializeApi(grpcClient, config)
	if err != nil {
		log.Fatalln(err)
	}

	server.Start(ginServer, config)
}
