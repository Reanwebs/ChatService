package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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

func HandleSocketConnection(c *gin.Context) {
	fmt.Println("entered to handle fun")
	userID := c.GetString("userId")
	fmt.Println(userID)

	if userID == "64ef152f4ca2c6fe73feaf9d" {
		socketID = "Edwin"
	} else {
		socketID = "EdwinV2"
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}

	connectedClients[socketID] = conn

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
			connectedMessage := []byte(wsMessage.Text)
			err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
			if err != nil {
				log.Println("Error sending message:", err)
				return
			}
			fmt.Println("Received message from", wsMessage.Sender)
			fmt.Println("Recipient:", wsMessage.Recipient)
			fmt.Println("Message text:", wsMessage.Text)
			continue
		}

	}
}
