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
	"github.com/seojoonrp/bbiyong-backend/apperr"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 모든 오리진에서 CORS 설정 허용
	},
}

type ChatHandler struct {
	hub            *ws.Hub
	chatService    services.ChatService
	userService    services.UserService
	meetingService services.MeetingService
}

func NewChatHandler(h *ws.Hub, cs services.ChatService, us services.UserService, ms services.MeetingService) *ChatHandler {
	return &ChatHandler{
		hub:            h,
		chatService:    cs,
		userService:    us,
		meetingService: ms,
	}
}

// 웹소켓 연결 진입점
func (h *ChatHandler) ChatConnect(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	meetingID := c.Param("id")

	if err := h.meetingService.VerifyParticipation(c.Request.Context(), meetingID, userID); err != nil {
		c.Error(err)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
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
		MeetingID:        meetingID,
		UserID:           userID,
		SenderName:       user.Nickname,
		SenderProfileURI: user.ProfileURI,
		ChatService:      h.chatService,
		Ctx:              wsCtx,
		Cancel:           cancel,
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	meetingID := c.Param("id")

	if err := h.meetingService.VerifyParticipation(c.Request.Context(), meetingID, userID); err != nil {
		c.Error(err)
		return
	}

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		c.Error(apperr.BadRequest("invalid limit parameter", err))
		return
	}

	// TODO : Last ID 등의 파라미터를 통해 페이징 처리 구현
	history, err := h.chatService.GetChatHistory(c.Request.Context(), meetingID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, history)
}
