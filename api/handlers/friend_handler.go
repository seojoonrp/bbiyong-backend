// api/handlers/friend_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/models"
)

type FriendHandler struct {
	friendService services.FriendService
}

func NewFriendHandler(fs services.FriendService) *FriendHandler {
	return &FriendHandler{friendService: fs}
}

func (h *FriendHandler) RequestFriend(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	targetID := c.Param("id")

	if err := h.friendService.RequestFriend(c.Request.Context(), userID, targetID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request successfully sent"})
}

func (h *FriendHandler) AcceptFriend(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	friendshipID := c.Param("id")

	if err := h.friendService.AcceptFriend(c.Request.Context(), userID, friendshipID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

func (h *FriendHandler) GetFriendList(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	status := c.Query("status")

	friends, err := h.friendService.ListFriends(c.Request.Context(), userID, status)
	if err != nil {
		c.Error(err)
		return
	}

	if friends == nil {
		friends = []models.FriendInfo{}
	}

	c.JSON(http.StatusOK, friends)
}
