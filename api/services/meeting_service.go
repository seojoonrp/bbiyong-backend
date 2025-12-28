// api/services/meeting_service.go

package services

import (
	"context"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MeetingService interface {
	CreateMeeting(ctx context.Context, hostID string, req models.CreateMeetingRequest) error
	GetNearbyMeetings(ctx context.Context, lon, lat float64, radius float64) ([]models.Meeting, error)
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
		Status:          "RECRUITING",
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
