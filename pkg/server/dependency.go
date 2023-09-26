package server

import (
	"chat/pkg/api/delivery"
	authclient "chat/pkg/api/delivery/authClient"
	websocket "chat/pkg/api/delivery/webSocket"
	"chat/pkg/api/repository"
	"chat/pkg/api/usecase"
	"log"
)

func InitializeApi(authClient authclient.AutharizationClientMethods, config Config) (*Server, error) {

	dbClient, err := ConnectPsqlDB(config)
	if err != nil {
		log.Println("error in connecting DB", err)
	}

	privateRepo := repository.NewPrivateChatRepo(dbClient)
	groupRepo := repository.NewGroupChatRepo(dbClient)

	privateUsecase := usecase.NewPrivateChatUsecase(privateRepo)
	groupUsecase := usecase.NewGroupChatUsecase(groupRepo)

	handler := delivery.NewChatHandler(authClient, privateUsecase, groupUsecase)
	wsHandler := websocket.NewWebSocketHandler(privateUsecase, groupUsecase)
	routes := delivery.NewChatRoutes(handler, wsHandler)

	server := NewServer(routes, config)

	return &server, nil
}
