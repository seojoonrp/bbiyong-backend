// api/handlers/friend_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendHandler struct {
	friendService services.FriendService
}

func NewFriendHandler(fs services.FriendService) *FriendHandler {
	return &FriendHandler{friendService: fs}
}

func (h *FriendHandler) RequestFriend(c *gin.Context) {
	myIDStr, _ := c.Get("user_id")
	targetIDStr := c.Param("id")

	myID, err := primitive.ObjectIDFromHex(myIDStr.(string))
	targetID, err := primitive.ObjectIDFromHex(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user or target id"})
		return
	}

	if err := h.friendService.RequestFriend(c.Request.Context(), myID, targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request successfully sent"})
}

func (h *FriendHandler) AcceptFriend(c *gin.Context) {
	myIDStr, _ := c.Get("user_id")
	friendshipIDStr := c.Param("id")

	myID, err := primitive.ObjectIDFromHex(myIDStr.(string))
	fID, err := primitive.ObjectIDFromHex(friendshipIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user or friendship id"})
		return
	}

	if err := h.friendService.AcceptFriend(c.Request.Context(), fID, myID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}
