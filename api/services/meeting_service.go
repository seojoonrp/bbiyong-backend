// api/services/meeting_service.go

package services

import (
	"context"
	"errors"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MeetingService interface {
	CreateMeeting(ctx context.Context, hostID string, req models.CreateMeetingRequest) error
	GetNearbyMeetings(ctx context.Context, lon, lat float64, radius float64) ([]models.Meeting, error)
	JoinMeeting(ctx context.Context, meetingID, userID string) error
	LeaveMeeting(ctx context.Context, meetingID, userID string) error
}

type meetingService struct {
	repo repositories.MeetingRepository
}

func NewMeetingService(repo repositories.MeetingRepository) MeetingService {
	return &meetingService{repo: repo}
}

func (s *meetingService) CreateMeeting(ctx context.Context, hostID string, req models.CreateMeetingRequest) error {
	hID, _ := primitive.ObjectIDFromHex(hostID)

	meeting := models.Meeting{
		Title:       req.Title,
		Description: req.Description,
		PlaceName:   req.PlaceName,
		Location: models.Location{
			Type:        "Point",
			Coordinates: []float64{req.Longitude, req.Latitude},
		},
		HostID:          hID,
		Participants:    []primitive.ObjectID{hID},
		MaxParticipants: req.MaxParticipants,
		AgeRange:        req.AgeRange,
		MeetingTime:     req.MeetingTime,
		Status:          models.MeetingStatusRecruiting,
		CreatedAt:       time.Now(),
	}

	return s.repo.Create(ctx, &meeting)
}

func (s *meetingService) GetNearbyMeetings(ctx context.Context, lon, lat float64, radius float64) ([]models.Meeting, error) {
	if radius == 0 {
		radius = 3000 // 기본 3km
	}
	return s.repo.FindNearby(ctx, lon, lat, radius)
}

func (s *meetingService) JoinMeeting(ctx context.Context, meetingID string, userID string) error {
	mID, _ := primitive.ObjectIDFromHex(meetingID)
	uID, _ := primitive.ObjectIDFromHex(userID)

	meeting, err := s.repo.FindByID(ctx, mID)
	if err != nil || meeting == nil {
		return errors.New("meeting not found")
	}

	success, err := s.repo.AddParticipant(ctx, mID, uID, meeting.MaxParticipants)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("failed to join the meeting")
	}

	return nil
}

func (s *meetingService) LeaveMeeting(ctx context.Context, meetingID string, userID string) error {
	mID, _ := primitive.ObjectIDFromHex(meetingID)
	uID, _ := primitive.ObjectIDFromHex(userID)

	meeting, err := s.repo.FindByID(ctx, mID)
	if err != nil || meeting == nil {
		return errors.New("meeting not found")
	}

	// 방장은 못나감
	if meeting.HostID == uID {
		return errors.New("host cannot leave the meeting")
	}

	success, err := s.repo.RemoveParticipant(ctx, mID, uID, meeting.MaxParticipants)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("failed to leave the meeting")
	}

	return nil
}
