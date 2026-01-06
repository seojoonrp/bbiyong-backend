// models/save_model.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Save struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	MeetingID primitive.ObjectID `bson:"meeting_id"`
	CreatedAt time.Time          `bson:"created_at"`
}
