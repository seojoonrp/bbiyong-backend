// api/handlers/chat_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/seojoonrp/bbiyong-backend/api/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 모든 오리진에서 CORS 설정 허용
	},
}

type ChatHandler struct {
	hub *ws.Hub
}

func NewChatHandler(h *ws.Hub) *ChatHandler {
	return &ChatHandler{hub: h}
}

// 웹소켓 연결 진입점
func (h *ChatHandler) ChatConnect(c *gin.Context) {
	meetingID := c.Param("id")
	userID, _ := c.Get("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Hub:       h.hub,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		MeetingID: meetingID,
		UserID:    userID.(string),
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
