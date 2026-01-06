// api/ws/client.go

package ws

import (
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/seojoonrp/bbiyong-backend/api/services"
)

type Client struct {
	Hub              *Hub
	Conn             *websocket.Conn
	Send             chan []byte
	MeetingID        string
	UserID           string
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

		savedMsg, err := c.ChatService.SaveMessage(context.Background(), c.MeetingID, c.UserID, string(message), c.SenderName, c.SenderProfileURI)
		if err != nil {
			c.sendError("Failed to save message:" + err.Error())
			continue
		}

		finalJson, err := json.Marshal(savedMsg)
		if err != nil {
			c.sendError("Internal data error:" + err.Error())
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

type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (c *Client) sendError(msg string) {
	errPayload, _ := json.Marshal(ErrorResponse{
		Type:    "ERROR",
		Message: msg,
	})
	c.Send <- errPayload
}
