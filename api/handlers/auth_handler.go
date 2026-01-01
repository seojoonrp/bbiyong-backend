// api/handlers/auth_handler.go

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isNewUser, token, user, err := h.authService.LoginWithGoogle(c.Request.Context(), req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to login with Google: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"user":        user,
		"isNewUser":   isNewUser,
	})
}

func (h *AuthHandler) KakaoLogin(c *gin.Context) {
	var input struct {
		AccessToken string `json:"accessToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isNew, token, user, err := h.authService.LoginWithKakao(c.Request.Context(), input.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kakao login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"user":        user,
		"isNewUser":   isNew,
	})
}

func (h *AuthHandler) AppleLogin(c *gin.Context) {
	var input struct {
		IdentityToken string `json:"identityToken" binding:"required"`
		FullName      struct {
			GivenName  string `json:"givenName"`
			FamilyName string `json:"familyName"`
		} `json:"fullName"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isNew, token, user, err := h.authService.LoginWithApple(c.Request.Context(), input.IdentityToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Apple login failed"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "username query parameter is required"})
		return
	}

	available, err := h.authService.IsUsernameAvailable(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username availability: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isAvailable": available})
}

func (h *AuthHandler) SetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.SetProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.CompleteProfile(c.Request.Context(), userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile completed successfully"})
}
