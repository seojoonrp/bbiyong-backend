// models/user_model.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ProviderLocal  = "LOCAL"
	ProviderKakao  = "KAKAO"
	ProviderGoogle = "GOOGLE"
	ProviderApple  = "APPLE"
)

const (
	GenderMale   = "MALE"
	GenderFemale = "FEMALE"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Password     string             `bson:"password,omitempty" json:"-"`
	Nickname     string             `bson:"nickname" json:"nickname"`
	ProfileURI   string             `bson:"profile_uri" json:"profileURI"`
	Age          int                `bson:"age" json:"age"`
	Gender       string             `bson:"gender" json:"gender"`
	Level        int                `bson:"level" json:"level"`
	Residence    string             `bson:"residence" json:"residence"`
	Provider     string             `bson:"provider" json:"provider"`
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
	Nickname   string `json:"nickname" binding:"required"`
	ProfileURI string `json:"profileURI" binding:"required"`
	Age        int    `json:"age" binding:"required"`
	Gender     string `json:"gender" binding:"required"`
	Residence  string `json:"residence" binding:"required"`
}
