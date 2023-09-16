package delivery

import (
	"chat/pkg/api/usecase"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	PrivateChatUsecase usecase.PrivateChatUsecaseMethods
	GroupChatUsecase   usecase.GroupChatUsecaseMethods
}

func NewChatHandler(privateUsecase usecase.PrivateChatUsecaseMethods, groupUsecase usecase.GroupChatUsecaseMethods) ChatHandler {
	return ChatHandler{
		PrivateChatUsecase: privateUsecase,
		GroupChatUsecase:   groupUsecase,
	}
}

type ChatHandlerMethods interface {
	GetPrivateChat(c *gin.Context)
	GroupChatStart(c *gin.Context)
}

func (h ChatHandler) GetPrivateChat(c *gin.Context) {

	h.PrivateChatUsecase.PrivateChatStart()
}

func (h ChatHandler) GetGroupChat(c *gin.Context) {
	h.GroupChatUsecase.GroupChatStart()
}
