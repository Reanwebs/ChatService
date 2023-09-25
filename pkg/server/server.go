package server

import (
	"chat/pkg/api/delivery"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Routes delivery.ChatRoutes
	Server *gin.Engine
}

func NewServer(routes delivery.ChatRoutes, config Config) Server {

	router := gin.Default()
	routes.SetPrivteChatRoutes(router)
	routes.SetGroupChatRoutes(router)
	routes.SetPublicChatRoutes(router)

	return Server{
		Routes: routes,
		Server: router,
	}
}

func Start(server *Server, config Config) {

	server.Server.Run(config.Port)

	log.Println("server started running on port", config.Port)
}
