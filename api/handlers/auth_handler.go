// api/handlers/auth_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/models"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.BadRequest("invalid request body", err))
		return
	}

	err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.BadRequest("invalid request body", err))
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"user":        user,
		"isNewUser":   false,
	})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req struct {
		IDToken string `json:"idToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.BadRequest("invalid request body", err))
		return
	}

	isNewUser, token, user, err := h.authService.LoginWithGoogle(c.Request.Context(), req.IDToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"user":        user,
		"isNewUser":   isNewUser,
	})
}

func (h *AuthHandler) KakaoLogin(c *gin.Context) {
	var req struct {
		AccessToken string `json:"accessToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.BadRequest("invalid request body", err))
		return
	}

	isNew, token, user, err := h.authService.LoginWithKakao(c.Request.Context(), req.AccessToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"user":        user,
		"isNewUser":   isNew,
	})
}

func (h *AuthHandler) AppleLogin(c *gin.Context) {
	var req struct {
		IdentityToken string `json:"identityToken" binding:"required"`
		FullName      struct {
			GivenName  string `json:"givenName"`
			FamilyName string `json:"familyName"`
		} `json:"fullName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.BadRequest("invalid request body", err))
		return
	}

	isNew, token, user, err := h.authService.LoginWithApple(c.Request.Context(), req.IdentityToken)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"user":        user,
		"isNewUser":   isNew,
	})
}

func (h *AuthHandler) CheckUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.Error(apperr.BadRequest("username query parameter is required", nil))
		return
	}

	available, err := h.authService.IsUsernameAvailable(c.Request.Context(), username)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"isAvailable": available})
}

func (h *AuthHandler) SetProfile(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req models.SetProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.BadRequest("invalid request body", err))
		return
	}

	err = h.authService.CompleteProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile completed successfully"})
}
