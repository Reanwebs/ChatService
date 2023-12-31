package delivery

import (
	"chat/pb/client"
	authclient "chat/pkg/api/delivery/authClient"
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/usecase"
	"errors"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	UserID string
	err    error
)

type ChatHandler struct {
	PrivateChatUsecase usecase.PrivateChatUsecaseMethods
	GroupChatUsecase   usecase.GroupChatUsecaseMethods
	AuthClient         authclient.AutharizationClientMethods
}

func NewChatHandler(authClient authclient.AutharizationClientMethods, privateUsecase usecase.PrivateChatUsecaseMethods, groupUsecase usecase.GroupChatUsecaseMethods) ChatHandlerMethods {
	return ChatHandler{
		PrivateChatUsecase: privateUsecase,
		GroupChatUsecase:   groupUsecase,
		AuthClient:         authClient,
	}
}

type ChatHandlerMethods interface {
	StartPrivateChat(c *gin.Context)
	GetPrivateChat(c *gin.Context)
	PrivateChatHistory(c *gin.Context)

	StartGroupChat(c *gin.Context)
	GetGroupChatList(c *gin.Context)
	GetGroupChatHistory(c *gin.Context)
}

func (h ChatHandler) StartPrivateChat(c *gin.Context) {
	userID := c.GetString("userId")
	body := models.StartChat{}
	if err = c.ShouldBindJSON(&body); err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, errors.Join(errors.New("JSON Binding failed"), err))
		return
	}
	res, err := h.AuthClient.GetUserDetails(c, &client.GetUserDetailsRequest{
		UserID: body.RecipientID,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errors.Join(errors.New("Auth client error"), err))
	}
	if res.Permission != true {
		c.JSON(http.StatusBadRequest, errors.Join(errors.New("User permission denied"), err))
	}
	input := models.PrivateChat{
		UserID:            userID,
		UserName:          body.UserName,
		RecipientID:       body.RecipientID,
		RecipientName:     res.UserName,
		RecipientAvatarID: res.AvatarID,
		StartAt:           time.Time{},
		LastSeen:          time.Time{},
	}
	if err = h.PrivateChatUsecase.StartChat(input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, "Private chat started")
}

func (h ChatHandler) GetPrivateChat(c *gin.Context) {
	input := models.GetChat{
		UserID: c.GetString("userId"),
	}
	response, err := h.PrivateChatUsecase.PrivateChatList(input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	sort.SliceStable(response, func(i, j int) bool {
		return response[i].LastSeen.After(response[j].LastSeen)
	})
	c.JSON(http.StatusOK, response)
}

func (h ChatHandler) PrivateChatHistory(c *gin.Context) {
	userID := c.GetString("userId")
	input := models.ChatHistory{}
	if err = c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, errors.Join(errors.New("JSON Binding failed"), err))
		return
	}
	sendedMessages, err := h.PrivateChatUsecase.RetrivePrivateChatHistory(userID, input.RecipientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.Join(errors.New("Server error"), err))
		return
	}
	recievedMessages, err := h.PrivateChatUsecase.RetriveRecievedChatHistory(userID, input.RecipientID)
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

// Group Chat Handlers
func (h ChatHandler) GetGroupChatList(c *gin.Context) {
	userID := c.GetString("userId")
	response, err := h.GroupChatUsecase.GetGroupList(userID)
	if err != nil {
		log.Println(err, "error retrieving grouplist")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	sort.SliceStable(response, func(i, j int) bool {
		return response[i].LastSeen.After(response[j].LastSeen)
	})
	c.JSON(http.StatusOK, response)
}

func (h ChatHandler) StartGroupChat(c *gin.Context) {
	userID := c.GetString("userId")
	var model models.GroupChat
	if err = c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errors.Join(errors.New("JSON Binding failed"), err))
		return
	}
	res, err := h.AuthClient.UserGroupPermission(c, &client.UserGroupPermissionRequest{
		UserID:  userID,
		GroupID: model.GroupID,
	})
	if res.Permission != true {
		c.JSON(http.StatusBadRequest, errors.Join(errors.New("Permission Denied"), err))
		return
	}
	input := models.GroupChat{
		UserID:        userID,
		UserName:      model.UserName,
		GroupID:       model.GroupID,
		GroupName:     res.GroupName,
		GroupAvatarID: res.GroupAvatarID,
		Permission:    res.Permission,
		StartAt:       time.Time{},
		LastSeen:      time.Time{},
	}
	if err = h.GroupChatUsecase.GroupChatStart(input); err != nil {
		log.Println(err, "error starting groupchat")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "Group chat started")
}

func (h ChatHandler) GetGroupChatHistory(c *gin.Context) {
	var model models.GroupChatHistory
	if err = c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errors.Join(errors.New("JSON Binding failed"), err))
		return
	}
	model.UserID = c.GetString("userId")
	response, err := h.GroupChatUsecase.GetGroupChatHistory(model)
	if err != nil {
		log.Println(err, "error retrieving groupchat")
		c.JSON(http.StatusBadRequest, err)
	}
	sort.SliceStable(response, func(i, j int) bool {
		return response[i].Time.Before(response[j].Time)
	})
	c.JSON(http.StatusOK, response)
}
