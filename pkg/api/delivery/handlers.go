package delivery

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/usecase"
	"errors"
	"log"
	"net/http"

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
	StartPrivateChat(c *gin.Context)
	GroupChatStart(c *gin.Context)
}

func (h ChatHandler) GetPrivateChat(c *gin.Context) {

	c.JSON(http.StatusAccepted, "Hi private chat route connected")
	h.PrivateChatUsecase.PrivateChatStart()
}

func (h ChatHandler) StartPrivateChat(c *gin.Context) {
	input := models.PrivateChat{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, errors.Join(errors.New("JSON Binding failed"), err))
		return
	}

	if err := h.PrivateChatUsecase.StartChat(input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, "Private chat started")
}

func (h ChatHandler) GetGroupChat(c *gin.Context) {
	h.GroupChatUsecase.GroupChatStart()
}
