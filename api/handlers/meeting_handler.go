// api/handlers/meeting_handler.go

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/models"
)

type MeetingHandler struct {
	service services.MeetingService
}

func NewMeetingHandler(s services.MeetingService) *MeetingHandler {
	return &MeetingHandler{service: s}
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateMeeting(c.Request.Context(), userID.(string), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "meeting successfully created"})
}

func (h *MeetingHandler) GetNearby(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("latitude"), 64)
	lon, _ := strconv.ParseFloat(c.Query("longitude"), 64)
	radius, _ := strconv.ParseFloat(c.Query("radius"), 64)

	meetings, err := h.service.GetNearbyMeetings(c.Request.Context(), lon, lat, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meetings)
}

func (h *MeetingHandler) Join(c *gin.Context) {
	meetingID := c.Param("id")
	userID, _ := c.Get("user_id")

	err := h.service.JoinMeeting(c.Request.Context(), meetingID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully joined the meeting"})
}

func (h *MeetingHandler) Leave(c *gin.Context) {
	meetingID := c.Param("id")
	userID, _ := c.Get("user_id")

	err := h.service.LeaveMeeting(c.Request.Context(), meetingID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully left the meeting"})
}
