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

	err = server.InitializeApi(config)
	if err != nil {
		log.Fatalln(err)
	}
}
