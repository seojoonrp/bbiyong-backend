// api/handlers/chat_handler.go

package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/api/ws"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 모든 오리진에서 CORS 설정 허용
	},
}

type ChatHandler struct {
	hub         *ws.Hub
	chatService services.ChatService
	userService services.UserService
}

func NewChatHandler(h *ws.Hub, cs services.ChatService, us services.UserService) *ChatHandler {
	return &ChatHandler{hub: h, chatService: cs, userService: us}
}

// 웹소켓 연결 진입점
func (h *ChatHandler) ChatConnect(c *gin.Context) {
	meetingID := c.Param("id")
	userID, _ := c.Get("user_id")

	mID, err := primitive.ObjectIDFromHex(meetingID)
	uID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	wsCtx, cancel := context.WithCancel(context.Background())

	client := &ws.Client{
		Hub:              h.hub,
		Conn:             conn,
		Send:             make(chan []byte, 256),
		MeetingID:        mID,
		UserID:           uID,
		SenderName:       user.Nickname,
		SenderProfileURI: "temp", // TODO
		ChatService:      h.chatService,
		Ctx:              wsCtx,
		Cancel:           cancel,
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
