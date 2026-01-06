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
	uIDVal, _ := c.Get("user_id")
	uIDStr, ok := uIDVal.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	mIDStr := c.Param("id")

	err := h.saveService.SaveMeeting(c.Request.Context(), uIDStr, mIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meeting saved successfully"})
}
