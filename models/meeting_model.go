// models/meeting_model.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Location GeoJSON 구조체
type Location struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

const (
	MeetingStatusRecruiting = "recruiting"
	MeetingStatusFull       = "full"
	MeetingStatusOngoing    = "ongoing"
	MeetingStatusFinished   = "finished"
)

type Meeting struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title           string               `bson:"title" json:"title"`
	Description     string               `bson:"description" json:"description"`
	PlaceName       string               `bson:"place_name" json:"placeName"` // 장소 이름 (ex. 율동공원, 서울대 총운동장)
	Location        Location             `bson:"location" json:"location"`
	HostID          primitive.ObjectID   `bson:"host_id" json:"hostID"`
	Participants    []primitive.ObjectID `bson:"participants" json:"participants"`
	MaxParticipants int                  `bson:"max_participants" json:"maxParticipants"`
	AgeRange        string               `bson:"age_range" json:"ageRange"`
	MeetingTime     time.Time            `bson:"meeting_time" json:"meetingTime"`
	Status          string               `bson:"status" json:"status"`
	CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}

type CreateMeetingRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	PlaceName       string    `json:"placeName" binding:"required"`
	Latitude        float64   `json:"latitude" binding:"required"`
	Longitude       float64   `json:"longitude" binding:"required"`
	MaxParticipants int       `json:"maxParticipants"`
	AgeRange        string    `json:"ageRange"`
	MeetingTime     time.Time `json:"meetingTime" binding:"required"`
}
