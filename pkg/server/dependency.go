package server

import (
	"log"
)

func InitializeApi(config Config) error {

	_, err := ConnectDB(config)
	if err != nil {
		log.Println("error in connecting DB", err)
	}

	// privateRepo := repository.NewPrivateChatRepo(DB)
	// groupRepo := repository.NewGroupChatRepo(DB)

	// privateUsecase := usecase.NewPrivateChatUsecase(privateRepo)

	// groupUsecase := usecase.NewGroupChatUsecase(groupRepo)

	// handler := delivery.NewChatHandler(privateUsecase, groupUsecase)

	// routes := delivery.NewChatRoutes(handler)

	return nil
}
