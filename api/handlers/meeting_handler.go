// api/handlers/meeting_handler.go

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/models"
)

type MeetingHandler struct {
	service services.MeetingService
}

func NewMeetingHandler(s services.MeetingService) *MeetingHandler {
	return &MeetingHandler{service: s}
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req models.CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.BadRequest("invalid request body", err))
		return
	}

	if err := h.service.CreateMeeting(c.Request.Context(), userID, req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "successfully created meeting"})
}

func (h *MeetingHandler) GetNearby(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("latitude"), 64)
	lon, err := strconv.ParseFloat(c.Query("longitude"), 64)
	radius, err := strconv.ParseFloat(c.Query("radius"), 64)
	if err != nil {
		c.Error(apperr.BadRequest("invalid location query parameters", err))
		return
	}

	daysStr := c.QueryArray("day_of_week")

	meetings, err := h.service.GetNearbyMeetings(c.Request.Context(), lon, lat, radius, daysStr)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, meetings)
}

func (h *MeetingHandler) Join(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	meetingID := c.Param("id")

	err = h.service.JoinMeeting(c.Request.Context(), meetingID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully joined the meeting"})
}

func (h *MeetingHandler) Leave(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	meetingID := c.Param("id")

	err = h.service.LeaveMeeting(c.Request.Context(), meetingID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully left the meeting"})
}
