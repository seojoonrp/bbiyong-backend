// api/handlers/friend_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/models"
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

func (h *FriendHandler) GetFriendList(c *gin.Context) {
	myIDStr, _ := c.Get("user_id")
	myID, err := primitive.ObjectIDFromHex(myIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	status := c.DefaultQuery("status", "ACCEPTED")

	friends, err := h.friendService.ListFriends(c.Request.Context(), myID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch friend list: " + err.Error()})
		return
	}

	if friends == nil {
		friends = []models.FriendInfo{}
	}

	c.JSON(http.StatusOK, friends)
}
