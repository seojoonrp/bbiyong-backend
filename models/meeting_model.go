// models/meeting_model.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MeetingStatusRecruiting = "RECRUITING"
	MeetingStatusFull       = "FULL"
	MeetingStatusOngoing    = "ONGOING"
	MeetingStatusFinished   = "FINISHED"
)

type Meeting struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title           string               `bson:"title" json:"title"`
	Description     string               `bson:"description" json:"description"`
	Category        string               `bson:"category" json:"category"`
	ImageURL        string               `bson:"image_url" json:"imageURL"`
	PlaceName       string               `bson:"place_name" json:"placeName"`
	Location        Location             `bson:"location" json:"location"`
	MeetingTime     time.Time            `bson:"meeting_time" json:"meetingTime"`
	DayOfWeek       int                  `bson:"day_of_week" json:"dayOfWeek"`
	AgeRange        [2]int               `bson:"age_range" json:"ageRange"`
	HostID          primitive.ObjectID   `bson:"host_id" json:"hostID"`
	Status          string               `bson:"status" json:"status"`
	ParticipantIDs  []primitive.ObjectID `bson:"participant_ids" json:"participantIDs"`
	MaxParticipants int                  `bson:"max_participants" json:"maxParticipants"`
	SaveCount       int                  `bson:"save_count" json:"saveCount"`
	CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}

type CreateMeetingRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Category        string    `json:"category" binding:"required"`
	ImageURL        string    `json:"imageURL" binding:"required"`
	PlaceName       string    `json:"placeName" binding:"required"`
	Location        Location  `json:"location" binding:"required"`
	MeetingTime     time.Time `json:"meetingTime" binding:"required"`
	DayOfWeek       int       `json:"dayOfWeek" binding:"required"`
	AgeRange        [2]int    `json:"ageRange" binding:"required"`
	MaxParticipants int       `json:"maxParticipants" binding:"required"`
}
