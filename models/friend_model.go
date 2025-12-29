// models/friend_model.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	FriendStatusPending  = "PENDING"
	FriendStatusAccepted = "ACCEPTED"
	FriendStatusRejected = "REJECTED"
)

type Friendship struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RequesterID primitive.ObjectID `bson:"requester_id" json:"requesterID"`
	AddresseeID primitive.ObjectID `bson:"addressee_id" json:"addresseeID"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}