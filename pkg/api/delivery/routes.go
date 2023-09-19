package delivery

import (
	websocket "chat/pkg/api/delivery/webSocket"

	"github.com/gin-gonic/gin"
)

type ChatRoutes struct {
	ChatHandler ChatHandler
	Middleware  MiddlewareMethods
}

func NewChatRoutes(handler ChatHandler) ChatRoutes {
	return ChatRoutes{
		ChatHandler: handler,
		Middleware:  Middleware{},
	}
}

func (h ChatRoutes) SetPrivteChatRoutes(router *gin.Engine) {

	router.Use(h.Middleware.AuthenticateUser)
	router.GET("ws", websocket.HandleSocketConnection)
	router.POST("chat/get", h.ChatHandler.GetPrivateChat)
	router.POST("chat/create", h.ChatHandler.StartPrivateChat)

}

func (h ChatRoutes) SetGroupChatRoutes(router *gin.Engine) {
	router.GET("/group/chat", h.ChatHandler.GetGroupChat)
}
