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
	AddPrivateChatHistory(string, string, string, string, *gin.Context)
}

func NewWebSocketHandler(privateUsecase usecase.PrivateChatUsecaseMethods, groupUsecase usecase.GroupChatUsecaseMethods) WebSocketMethods {
	return WebSocketHandler{
		PrivateChatUsecase: privateUsecase,
		GroupChatUsecase:   groupUsecase,
	}
}

var (
	connectedClients      = make(map[string]*websocket.Conn)
	connectedGroupClients = make(map[string]*websocket.Conn)
	socketID              string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketMessage struct {
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
	Online    bool      `json:"online"`
}

func (w WebSocketHandler) HandleSocketConnection(c *gin.Context) {
	userName := c.DefaultQuery("userName", "")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}
	defer func() {
		delete(connectedClients, userName)
		conn.Close()
	}()
	if _, ok := connectedClients[userName]; !ok {
		connectedMessage := []byte("Connected to the server")
		err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		connectedClients[userName] = conn
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
				if err = sendOnlineStatus(true, userName); err != nil {
					log.Println("Error forwarding onlineStatus to user:", err)
				}
				w.AddPrivateChatHistory(userName, recipient, "delivered", wsMessage.Text, c)
			} else {
				log.Println("Recipient is not connected")
				if err = sendOnlineStatus(false, userName); err != nil {
					log.Println("Error forwarding onlineStatus to user:", err)
				}
				w.AddPrivateChatHistory(userName, recipient, "undelivered", wsMessage.Text, c)
			}
		}
	}
}

func (w WebSocketHandler) HandleGroupSocketConnection(c *gin.Context) {
	userName := c.DefaultQuery("userName", "")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}
	defer func() {
		delete(connectedGroupClients, userName)
		conn.Close()
	}()
	if _, ok := connectedGroupClients[userName]; !ok {
		connectedMessage := []byte("Connected to the server")
		err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		connectedGroupClients[userName] = conn
	}
}

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

func (w WebSocketHandler) AddPrivateChatHistory(userId string, recipientId string, status string, Text string, c *gin.Context) {
	input := models.PrivateChatHistory{
		Text:   Text,
		Status: status,
		Time:   time.Now(),
	}
	if err := w.PrivateChatUsecase.CreatePrivateChatHistory(userId, recipientId, input); err != nil {
		log.Println(err)
		fmt.Println(err)
	}
}
