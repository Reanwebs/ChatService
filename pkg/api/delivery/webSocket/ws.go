package websocket

import (
	"chat/pkg/api/delivery/models"
	"chat/pkg/api/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	WebSocketMethods
	PrivateChatUsecase usecase.PrivateChatUsecaseMethods
	GroupChatUsecase   usecase.GroupChatUsecaseMethods
}

type WebSocketMethods interface {
	HandleSocketConnection(*gin.Context)
	HandleGroupSocketConnection(c *gin.Context)
	HandlePublicSocketConnection(c *gin.Context)
	AddPrivateChatHistory(string, string, string, string, string, *gin.Context)
	AddGroupChatHistory(string, string, models.WebSocketGroupMessage, *gin.Context)
}

func NewWebSocketHandler(privateUsecase usecase.PrivateChatUsecaseMethods, groupUsecase usecase.GroupChatUsecaseMethods) WebSocketMethods {
	return WebSocketHandler{
		PrivateChatUsecase: privateUsecase,
		GroupChatUsecase:   groupUsecase,
	}
}

var (
	connectedClients       = make(map[string]*websocket.Conn)
	connectedGroupClients  = make(map[string]map[string]*websocket.Conn)
	connectedPublicClients = make(map[string]*websocket.Conn)
	socketID               string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w WebSocketHandler) HandleSocketConnection(c *gin.Context) {
	userID := c.GetString("userId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}
	defer func() {
		delete(connectedClients, userID)
		conn.Close()
	}()
	if _, ok := connectedClients[userID]; !ok {
		connectedMessage := []byte("Connected to the server")
		err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		connectedClients[userID] = conn
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err, "socket error")
			return
		}
		if len(p) == 0 {
			log.Println("Empty WebSocket message received")
			continue
		}
		var wsMessage models.WebSocketMessage
		if err := json.Unmarshal(p, &wsMessage); err != nil {
			log.Println("Error decoding WebSocket message:", err)
		} else {
			recipient := wsMessage.Recipient
			if recipientConn, ok := connectedClients[recipient]; ok {
				err := recipientConn.WriteMessage(websocket.TextMessage, p)
				if err != nil {
					log.Println("Error forwarding message to recipient:", err)
				}
				if err = sendOnlineStatus(true, userID); err != nil {
					log.Println("Error forwarding onlineStatus to user:", err)
				}
				w.AddPrivateChatHistory(wsMessage.User, userID, recipient, "delivered", wsMessage.Text, c)
			} else {
				log.Println("Recipient is not connected")
				if err = sendOnlineStatus(false, userID); err != nil {
					log.Println("Error forwarding onlineStatus to user:", err)
				}
				w.AddPrivateChatHistory(wsMessage.User, userID, recipient, "undelivered", wsMessage.Text, c)
			}
		}
	}
}

func (w WebSocketHandler) HandleGroupSocketConnection(c *gin.Context) {
	groupID := c.Query("groupName")
	userID := c.GetString("userId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Socket connection error:", err)
		return
	}
	defer func() {
		delete(connectedGroupClients[groupID], userID)
		conn.Close()
	}()
	if connectedGroupClients[groupID] == nil {
		connectedGroupClients[groupID] = make(map[string]*websocket.Conn)
	}

	connectedGroupClients[groupID][userID] = conn

	if _, ok := connectedGroupClients[groupID][userID]; !ok {
		connectedMessage := []byte("Connected to the group chat")
		err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Socket error:", err)
			return
		}
		if len(p) == 0 {
			log.Println("Empty WebSocket message received")
			continue
		}
		var wsMessage models.WebSocketGroupMessage
		if err := json.Unmarshal(p, &wsMessage); err != nil {
			log.Println("Error decoding WebSocket message:", err)
		}
		for clientID, clientConn := range connectedGroupClients[wsMessage.GroupID] {
			if clientID != userID {
				err := clientConn.WriteMessage(websocket.TextMessage, p)
				if err != nil {
					log.Println("Error forwarding message to group member:", err)
				}
			}
		}
		w.AddGroupChatHistory(userID, groupID, wsMessage, c)
	}
}

func (w WebSocketHandler) HandlePublicSocketConnection(c *gin.Context) {
	userId := c.GetString("userId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}
	defer func() {
		delete(connectedPublicClients, userId)
		conn.Close()
	}()
	if _, ok := connectedPublicClients[userId]; !ok {
		connectedMessage := []byte("Connected to the server")
		err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		connectedPublicClients[userId] = conn
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		fmt.Println(p)
		for clientName, clientConn := range connectedPublicClients {
			if clientName != userId {
				err = clientConn.WriteMessage(messageType, p)
				if err != nil {
					log.Println("Error forwarding message:", err)
					return
				}
			}
		}
	}
}

// Non Ws Functions

func sendOnlineStatus(status bool, user string) error {
	onlineStatusMessage := models.WebSocketMessage{
		Type:      "onlineStatus",
		Sender:    user,
		Recipient: "",
		Text:      "",
		Time:      time.Time{},
		Online:    status,
	}
	message, err := json.Marshal(onlineStatusMessage)
	if err != nil {
		log.Println("Error marshaling online status:", err)
		return err
	}

	userConn, _ := connectedClients[user]
	err = userConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Error sending online status:", err)
		return err
	}
	return nil
}

func (w WebSocketHandler) AddPrivateChatHistory(userName string, userId string, recipientId string, status string, Text string, c *gin.Context) {
	input := models.PrivateChatHistory{
		UserName:    userName,
		UserID:      userId,
		RecipientID: recipientId,
		Text:        Text,
		Status:      status,
		Time:        time.Now(),
	}
	if err := w.PrivateChatUsecase.CreatePrivateChatHistory(input); err != nil {
		log.Println(err)
		fmt.Println(err)
	}
}

func (w WebSocketHandler) AddGroupChatHistory(userId string, groupId string, message models.WebSocketGroupMessage, c *gin.Context) {
	input := models.GroupChatHistory{
		UserID:    userId,
		UserName:  message.SenderName,
		GroupID:   groupId,
		GroupName: message.GroupName,
		Text:      message.Text,
		Status:    "delivered",
		Time:      time.Time{},
	}
	if err := w.GroupChatUsecase.AddGroupChatHistory(input); err != nil {
		log.Println(err)
		fmt.Println(err)
	}
}
