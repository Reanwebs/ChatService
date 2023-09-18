package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketMessage struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

func HandleSocketConnection(c *gin.Context) {
	fmt.Println("entered to handle fun")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
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
			continue
		}
		fmt.Println("Received message from", wsMessage.Sender)
		fmt.Println("Recipient:", wsMessage.Recipient)
		fmt.Println("Message text:", wsMessage.Text)
	}
}
