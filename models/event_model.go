// models/event_model.go

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	EventJoinMeeting  = "JOIN"
	EventLeaveMeeting = "LEAVE"
)

type MeetingEvent struct {
	Type      string
	MeetingID primitive.ObjectID
	UserID    primitive.ObjectID
}
