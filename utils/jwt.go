// utils/jwt.go

package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seojoonrp/bbiyong-backend/config"
)

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
