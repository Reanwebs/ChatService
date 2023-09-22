package delivery

import (
	websocket "chat/pkg/api/delivery/webSocket"

	"github.com/gin-gonic/gin"
)

type ChatRoutes struct {
	ChatHandler      ChatHandlerMethods
	WebSocketHandler websocket.WebSocketMethods
	Middleware       MiddlewareMethods
}

func NewChatRoutes(handler ChatHandlerMethods, wsHandler websocket.WebSocketMethods) ChatRoutes {
	return ChatRoutes{
		ChatHandler:      handler,
		WebSocketHandler: wsHandler,
		Middleware:       Middleware{},
	}
}

func (h ChatRoutes) SetPrivteChatRoutes(router *gin.Engine) {
	router.Use(h.Middleware.AuthenticateUser)
	router.GET("ws", h.WebSocketHandler.HandleSocketConnection)
	router.POST("chat/get-chatlist", h.ChatHandler.GetPrivateChat)
	router.POST("chat/create-chat", h.ChatHandler.StartPrivateChat)
	router.POST("chat/get-chat", h.ChatHandler.PrivateChatHistory)

}

func (h ChatRoutes) SetGroupChatRoutes(router *gin.Engine) {
	router.Use(h.Middleware.AuthenticateUser)
	router.GET("ws/group", h.WebSocketHandler.HandleGroupSocketConnection)
	router.POST("chat/group/start", h.ChatHandler.StartGroupChat)
	router.POST("chat/get-group", h.ChatHandler.GetGroupChatList)
	router.POST("chat/group/get-chat", h.ChatHandler.GetGroupChatHistory)
}
