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
	AddPrivateChatHistory(string, string, string, string, *gin.Context)
}

func NewWebSocketHandler(privateUsecase usecase.PrivateChatUsecaseMethods, groupUsecase usecase.GroupChatUsecaseMethods) WebSocketMethods {
	return WebSocketHandler{
		PrivateChatUsecase: privateUsecase,
		GroupChatUsecase:   groupUsecase,
	}
}

var (
	connectedClients = make(map[string]*websocket.Conn)
	socketID         string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketMessage struct {
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
}

func (w WebSocketHandler) HandleSocketConnection(c *gin.Context) {
	userName := c.DefaultQuery("userName", "")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}

	connectedClients[userName] = conn

	connectedMessage := []byte("Connected to the server")
	err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
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
				w.AddPrivateChatHistory(userName, recipient, "delivered", wsMessage.Text, c)
			} else {
				log.Println("Recipient is not connected")

				w.AddPrivateChatHistory(userName, recipient, "undelivered", wsMessage.Text, c)
			}
		}
	}
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
