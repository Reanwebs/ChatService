package delivery

import (
	"github.com/gin-gonic/gin"
)

type ChatRoutes struct {
	ChatHandler ChatHandler
}

func NewChatRoutes(handler ChatHandler) ChatRoutes {
	return ChatRoutes{
		ChatHandler: handler,
	}
}

func (h ChatRoutes) SetPrivteChatRoutes(router *gin.Engine) {

	router.GET("/private/chat", h.ChatHandler.GetPrivateChat)

}

func (h ChatRoutes) SetGroupChatRoutes(router *gin.Engine) {
	router.GET("/group/chat", h.ChatHandler.GetGroupChat)
}
