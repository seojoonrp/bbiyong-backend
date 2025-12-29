// api/services/chat_service.go

package services

import (
	"context"
	"errors"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService interface {
	SaveMessage(ctx context.Context, mID, uID primitive.ObjectID, content, name, profile string) (*models.ChatMessage, error)
	GetChatHistory(ctx context.Context, mID, uID primitive.ObjectID, limit int64) ([]models.ChatMessage, error)
}

type chatService struct {
	chatRepo    repositories.ChatRepository
	userRepo    repositories.UserRepository
	meetingRepo repositories.MeetingRepository
}

func NewChatService(cr repositories.ChatRepository, ur repositories.UserRepository, mr repositories.MeetingRepository) ChatService {
	return &chatService{chatRepo: cr, userRepo: ur, meetingRepo: mr}
}

func (s *chatService) SaveMessage(ctx context.Context, mID, uID primitive.ObjectID, content, name, profile string) (*models.ChatMessage, error) {
	msg := &models.ChatMessage{
		ID:               primitive.NewObjectID(),
		MeetingID:        mID,
		SenderID:         uID,
		SenderName:       name,
		SenderProfileURI: profile,
		Content:          content,
		Type:             models.ChatTypeTalk,
		CreatedAt:        time.Now(),
	}

	err := s.chatRepo.SaveMessage(ctx, msg)
	return msg, err
}

func (s *chatService) GetChatHistory(ctx context.Context, mID, uID primitive.ObjectID, limit int64) ([]models.ChatMessage, error) {
	meeting, err := s.meetingRepo.FindByID(ctx, mID)
	if err != nil {
		return nil, err
	}
	if meeting == nil {
		return nil, errors.New("meeting not found")
	}

	isParticipant := false
	for _, pID := range meeting.Participants {
		if pID == uID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, errors.New("user is not a participant of the meeting")
	}

	return s.chatRepo.GetChatHistory(ctx, mID, limit)
}
