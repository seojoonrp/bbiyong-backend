// api/handlers/chat_handler.go

package handlers

import (
	"context"
	"net/http"
	"strconv"

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

	isParticipant, err := h.chatService.CheckParticipation(c.Request.Context(), mID, uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check participation"})
		return
	}
	if !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not a participant of the meeting"})
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

func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	meetingIDStr := c.Param("id")
	userID, _ := c.Get("user_id")

	mID, err := primitive.ObjectIDFromHex(meetingIDStr)
	uID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 50
	}

	// TODO : Last ID 등의 파라미터를 통해 페이징 처리 구현
	history, err := h.chatService.GetChatHistory(c.Request.Context(), mID, uID, limit)
	if err != nil {
		if err.Error() == "you are not a participant of the meeting" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch chat history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
