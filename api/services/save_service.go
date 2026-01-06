// api/services/save_service.go

package services

import (
	"context"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SaveService interface {
	SaveMeeting(ctx context.Context, userID, meetingID string) error
}

type saveService struct {
	saveRepo    repositories.SaveRepository
	meetingRepo repositories.MeetingRepository
}

func NewSaveService(sr repositories.SaveRepository, mr repositories.MeetingRepository) SaveService {
	return &saveService{
		saveRepo:    sr,
		meetingRepo: mr,
	}
}

func (s *saveService) SaveMeeting(ctx context.Context, userID, meetingID string) error {
	uID, err := primitive.ObjectIDFromHex(userID)
	mID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		return apperr.BadRequest("invalid ID format", err)
	}

	err = s.saveRepo.Create(ctx, &models.Save{
		UserID:    uID,
		MeetingID: mID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return apperr.Conflict("meeting already saved", err)
		}
		return apperr.InternalServerError("failed to save meeting", err)
	}

	err = s.meetingRepo.IncrementSaveCount(ctx, mID)
	if err != nil {
		return apperr.InternalServerError("failed to increment save count", err)
	}

	return nil
}
