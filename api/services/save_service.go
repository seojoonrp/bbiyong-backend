// api/services/save_service.go

package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
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
		return errors.New("invalid ID format")
	}

	err = s.saveRepo.Create(ctx, &models.Save{
		UserID:    uID,
		MeetingID: mID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("meeting already saved")
		}
		return err
	}

	err = s.meetingRepo.IncrementSaveCount(ctx, mID)
	if err != nil {
		return fmt.Errorf("failed to increment save count: %v", err)
	}

	return nil
}
