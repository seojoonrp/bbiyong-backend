// models/event_model.go

package models

const (
	EventJoinMeeting  = "JOIN"
	EventLeaveMeeting = "LEAVE"
)

type MeetingEvent struct {
	Type      string
	MeetingID string
	UserID    string
}
