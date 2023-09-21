package delivery

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/usecase"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

var UserID string

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
	PrivateChatHistory(c *gin.Context)
	GroupChatStart(c *gin.Context)
	AddPrivateChatHistory(c *gin.Context)
}

func (h ChatHandler) GetPrivateChat(c *gin.Context) {
	input := models.GetChat{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, errors.Join(errors.New("JSON Binding failed"), err))
	}
	response, err := h.PrivateChatUsecase.PrivateChatList(input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Set("userName", input.UserID)
	c.JSON(http.StatusOK, response)
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

func (h ChatHandler) PrivateChatHistory(c *gin.Context) {
	input := models.PrivateChat{}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, errors.Join(errors.New("JSON Binding failed"), err))
		return
	}
	fmt.Println(input.RecipientID, "recipient for chat history")
	sendedMessages, err := h.PrivateChatUsecase.RetrivePrivateChatHistory(input.UserID, input.RecipientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.Join(errors.New("Server error"), err))
		return
	}
	recievedMessages, err := h.PrivateChatUsecase.RetriveRecievedChatHistory(input.UserID, input.RecipientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.Join(errors.New("Server error"), err))
		return
	}
	allMessages := append(sendedMessages, recievedMessages...)
	sort.SliceStable(allMessages, func(i, j int) bool {
		return allMessages[i].Time.Before(allMessages[j].Time)
	})
	response := models.ChatHistoryResponse{
		Messages: allMessages,
	}
	c.JSON(http.StatusOK, response)
}

//Group Chat Handlers

func (h ChatHandler) GetGroupChat(c *gin.Context) {
	h.GroupChatUsecase.GroupChatStart()
}
