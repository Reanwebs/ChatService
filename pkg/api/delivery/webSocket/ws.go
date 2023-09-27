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
	AddGroupChatHistory(string, string, WebSocketGroupMessage, *gin.Context)
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

type WebSocketMessage struct {
	User      string    `json:"user"`
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
	Online    bool      `json:"online"`
}

type WebSocketGroupMessage struct {
	Text       string `json:"text"`
	SenderName string `json:"sender"`
	GroupName  string `json:"recipient"`
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
		var wsMessage WebSocketMessage
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
	groupName := "Golang"
	userName := c.GetString("userId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Socket connection error:", err)
		return
	}
	defer func() {
		delete(connectedGroupClients[groupName], userName)
		conn.Close()
	}()
	if connectedGroupClients[groupName] == nil {
		connectedGroupClients[groupName] = make(map[string]*websocket.Conn)
	}

	connectedGroupClients[groupName][userName] = conn

	if _, ok := connectedGroupClients[groupName][userName]; !ok {
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

		for clientName, clientConn := range connectedGroupClients[groupName] {
			if clientName != userName {
				err := clientConn.WriteMessage(websocket.TextMessage, p)
				if err != nil {
					log.Println("Error forwarding message to group member:", err)
				}
			}
		}
		var wsMessage WebSocketGroupMessage
		if err := json.Unmarshal(p, &wsMessage); err != nil {
			log.Println("Error decoding WebSocket message:", err)
		}

		w.AddGroupChatHistory(userName, groupName, wsMessage, c)
	}
}

func (w WebSocketHandler) HandlePublicSocketConnection(c *gin.Context) {
	userName := c.DefaultQuery("userName", "")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}
	defer func() {
		delete(connectedPublicClients, userName)
		conn.Close()
	}()
	if _, ok := connectedPublicClients[userName]; !ok {
		connectedMessage := []byte("Connected to the server")
		err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		connectedPublicClients[userName] = conn
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		message := string(p)
		for clientName, clientConn := range connectedPublicClients {
			if clientName != userName {
				err = clientConn.WriteMessage(messageType, []byte(userName+": "+message))
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
	onlineStatusMessage := WebSocketMessage{
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

func (w WebSocketHandler) AddGroupChatHistory(userId string, groupId string, message WebSocketGroupMessage, c *gin.Context) {
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
