// models/chat_model.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ChatTypeTalk  = "TALK"
	ChatTypeJoin  = "JOIN"
	ChatTypeLeave = "LEAVE"
)

type ChatMessage struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MeetingID        primitive.ObjectID `bson:"meeting_id" json:"meetingID"`
	SenderID         primitive.ObjectID `bson:"sender_id" json:"senderID"`
	SenderName       string             `bson:"sender_name" json:"senderName"`
	SenderProfileURI string             `bson:"sender_profile_uri" json:"senderProfileUri"`
	Content          string             `bson:"content" json:"content"`
	Type             string             `bson:"type" json:"type"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
}
