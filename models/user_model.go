// models/user_model.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Password     string             `bson:"password" json:"-"`
	Nickname     string             `bson:"nickname" json:"nickname"`
	Provider     string             `bson:"provider" json:"provider"` // local, kakao, google, apple
	SocialID     string             `bson:"social_id,omitempty" json:"socialID,omitempty"`
	SocialEmail  string             `bson:"social_email,omitempty" json:"socialEmail,omitempty"`
	IsProfileSet bool               `bson:"is_profile_set" json:"isProfileSet"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SetProfileRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	// TODO : 나중에 프사, 생년월일, 성별 등 추가해야됨
}
