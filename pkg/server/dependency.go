package server

import (
	"chat/pkg/api/delivery"
	"chat/pkg/api/repository"
	"chat/pkg/api/usecase"
	"log"
)

func InitializeApi(config Config) (*Server, error) {

	dbClient, err := ConnectDB(config)
	if err != nil {
		log.Println("error in connecting DB", err)
	}

	privateRepo := repository.NewPrivateChatRepo(dbClient)
	groupRepo := repository.NewGroupChatRepo(dbClient)

	privateUsecase := usecase.NewPrivateChatUsecase(privateRepo)
	groupUsecase := usecase.NewGroupChatUsecase(groupRepo)

	handler := delivery.NewChatHandler(privateUsecase, groupUsecase)

	routes := delivery.NewChatRoutes(handler)

	server := NewServer(routes, config)

	return &server, nil
}
