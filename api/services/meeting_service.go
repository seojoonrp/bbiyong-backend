// api/services/meeting_service.go

package services

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MeetingService interface {
	CreateMeeting(ctx context.Context, hostID string, req models.CreateMeetingRequest) error
	GetNearbyMeetings(ctx context.Context, lon, lat float64, radius float64, days []string) ([]models.Meeting, error)
	VerifyParticipation(ctx context.Context, meetingID, userID string) error
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
		return apperr.BadRequest("invalid ID format", err)
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

	err = s.meetingRepo.Create(ctx, &meeting)
	if err != nil {
		return apperr.InternalServerError("failed to create meeting", err)
	}

	return nil
}

func (s *meetingService) GetNearbyMeetings(ctx context.Context, lon, lat float64, radius float64, days []string) ([]models.Meeting, error) {
	if radius == 0 {
		log.Println("Radius not provided, defaulting to 3000 meters")
		radius = 3000 // 기본 3km
	}

	var daysInt []int
	for _, s := range days {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, apperr.BadRequest("invalid day of week format", err)
		}

		if val >= 0 && val <= 6 {
			daysInt = append(daysInt, val)
		} else {
			return nil, apperr.BadRequest("invalid day of week value", nil)
		}
	}

	return s.meetingRepo.FindNearby(ctx, lon, lat, radius, daysInt)
}

func (s *meetingService) VerifyParticipation(ctx context.Context, meetingID, userID string) error {
	mID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		return apperr.BadRequest("invalid meeting ID format", err)
	}

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperr.InternalServerError("invalid user ID in token", err)
	}

	meeting, err := s.meetingRepo.FindByID(ctx, mID)
	if err != nil {
		return apperr.InternalServerError("failed to fetch meeting", err)
	}
	if meeting == nil {
		return apperr.NotFound("meeting not found", nil)
	}

	for _, pID := range meeting.ParticipantIDs {
		if pID == uID {
			return nil
		}
	}

	return apperr.Forbidden("you are not a participant of the meeting", nil)
}

func (s *meetingService) JoinMeeting(ctx context.Context, meetingID, userID string) error {
	mID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		return apperr.BadRequest("invalid meeting ID format", err)
	}

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperr.InternalServerError("invalid user ID in token", err)
	}

	meeting, err := s.meetingRepo.FindByID(ctx, mID)
	if err != nil {
		return apperr.InternalServerError("failed to fetch meeting", err)
	}
	if meeting == nil {
		return apperr.NotFound("meeting not found", nil)
	}

	success, err := s.meetingRepo.AddParticipant(ctx, mID, uID, meeting.MaxParticipants)
	if err != nil {
		return apperr.InternalServerError("failed to add participant", err)
	}
	if !success {
		return apperr.BadRequest("failed to join the meeting", errors.New("meeting may be full or user already joined"))
	}

	s.eventChan <- models.MeetingEvent{
		Type:      models.EventJoinMeeting,
		MeetingID: meetingID,
		UserID:    userID,
	}

	return nil
}

func (s *meetingService) LeaveMeeting(ctx context.Context, meetingID, userID string) error {
	mID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		return apperr.BadRequest("invalid meeting ID format", err)
	}

	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperr.InternalServerError("invalid user ID in token", err)
	}

	meeting, err := s.meetingRepo.FindByID(ctx, mID)
	if err != nil {
		return apperr.InternalServerError("failed to fetch meeting", err)
	}
	if meeting == nil {
		return apperr.NotFound("meeting not found", nil)
	}

	// 방장은 못나감
	if meeting.HostID == uID {
		return apperr.BadRequest("host cannot leave the meeting", nil)
	}

	success, err := s.meetingRepo.RemoveParticipant(ctx, mID, uID, meeting.MaxParticipants)
	if err != nil {
		return apperr.InternalServerError("failed to remove participant", err)
	}
	if !success {
		return apperr.BadRequest("failed to leave the meeting", errors.New("user may not be a participant"))
	}

	return nil
}
