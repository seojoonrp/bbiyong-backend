// api/ws/client.go

package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/seojoonrp/bbiyong-backend/api/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	Hub              *Hub
	Conn             *websocket.Conn
	Send             chan []byte
	MeetingID        primitive.ObjectID
	UserID           primitive.ObjectID
	SenderName       string
	SenderProfileURI string
	ChatService      services.ChatService
	Ctx              context.Context
	Cancel           context.CancelFunc
}

// 메시지를 읽어서 허브로 보냄
func (c *Client) ReadPump() {
	defer func() {
		c.Cancel()
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break // 연결이 뭔가 이상해졌을 때
		}

		savedMsg, err := c.ChatService.SaveMessage(c.Ctx, c.MeetingID, c.UserID, string(message), c.SenderName, c.SenderProfileURI)
		if err != nil {
			log.Println("Error saving message:", err)
			continue
		}

		finalJson, err := json.Marshal(savedMsg)
		if err != nil {
			log.Println("Error marshaling message:", err)
			continue
		}

		c.Hub.Broadcast <- MessagePayload{
			MeetingID: c.MeetingID,
			Data:      finalJson,
		}
	}
}

// 허브로부터 받은 메시지를 전송
func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
