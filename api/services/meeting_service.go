// api/services/meeting_service.go

package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MeetingService interface {
	CreateMeeting(ctx context.Context, hostID string, req models.CreateMeetingRequest) error
	GetNearbyMeetings(ctx context.Context, lon, lat float64, radius float64, days []string) ([]models.Meeting, error)
	JoinMeeting(ctx context.Context, meetingID, userID string) error
	LeaveMeeting(ctx context.Context, meetingID, userID string) error
}

type meetingService struct {
	meetingRepo repositories.MeetingRepository
	eventChan   chan<- models.MeetingEvent
}

func NewMeetingService(repo repositories.MeetingRepository, ec chan<- models.MeetingEvent) MeetingService {
	return &meetingService{meetingRepo: repo, eventChan: ec}
}

func (s *meetingService) CreateMeeting(ctx context.Context, hostID string, req models.CreateMeetingRequest) error {
	hID, err := primitive.ObjectIDFromHex(hostID)
	if err != nil {
		return errors.New("invalid host ID format")
	}

	meeting := models.Meeting{
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		ImageURL:        req.ImageURL,
		PlaceName:       req.PlaceName,
		Location:        req.Location,
		MeetingTime:     req.MeetingTime,
		DayOfWeek:       req.DayOfWeek,
		AgeRange:        req.AgeRange,
		HostID:          hID,
		Status:          models.MeetingStatusRecruiting,
		ParticipantIDs:  []primitive.ObjectID{hID},
		MaxParticipants: req.MaxParticipants,
		SaveCount:       0,
		CreatedAt:       time.Now(),
	}

	return s.meetingRepo.Create(ctx, &meeting)
}

func (s *meetingService) GetNearbyMeetings(ctx context.Context, lon, lat float64, radius float64, days []string) ([]models.Meeting, error) {
	if radius == 0 {
		radius = 3000 // 기본 3km
	}

	var daysInt []int
	for _, s := range days {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.New("invalid day of week format")
		}

		if val >= 0 && val <= 6 {
			daysInt = append(daysInt, val)
		} else {
			return nil, errors.New("invalid day of week value")
		}
	}

	return s.meetingRepo.FindNearby(ctx, lon, lat, radius, daysInt)
}

func (s *meetingService) JoinMeeting(ctx context.Context, meetingID string, userID string) error {
	mID, err := primitive.ObjectIDFromHex(meetingID)
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid ID format")
	}

	meeting, err := s.meetingRepo.FindByID(ctx, mID)
	if err != nil || meeting == nil {
		return errors.New("meeting not found")
	}

	success, err := s.meetingRepo.AddParticipant(ctx, mID, uID, meeting.MaxParticipants)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("failed to join the meeting")
	}

	s.eventChan <- models.MeetingEvent{
		Type:      models.EventJoinMeeting,
		MeetingID: mID,
		UserID:    uID,
	}

	return nil
}

func (s *meetingService) LeaveMeeting(ctx context.Context, meetingID string, userID string) error {
	mID, _ := primitive.ObjectIDFromHex(meetingID)
	uID, _ := primitive.ObjectIDFromHex(userID)

	meeting, err := s.meetingRepo.FindByID(ctx, mID)
	if err != nil || meeting == nil {
		return errors.New("meeting not found")
	}

	// 방장은 못나감
	if meeting.HostID == uID {
		return errors.New("host cannot leave the meeting")
	}

	success, err := s.meetingRepo.RemoveParticipant(ctx, mID, uID, meeting.MaxParticipants)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("failed to leave the meeting")
	}

	return nil
}
