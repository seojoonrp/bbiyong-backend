// api/handlers/save_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
)

type SaveHandler struct {
	saveService services.SaveService
}

func NewSaveHandler(ss services.SaveService) *SaveHandler {
	return &SaveHandler{saveService: ss}
}

func (h *SaveHandler) SaveMeeting(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	mIDStr := c.Param("id")

	err = h.saveService.SaveMeeting(c.Request.Context(), userID, mIDStr)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meeting saved successfully"})
}

func (h *SaveHandler) UnsaveMeeting(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	mIDStr := c.Param("id")

	err = h.saveService.UnsaveMeeting(c.Request.Context(), userID, mIDStr)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meeting unsaved successfully"})
}
