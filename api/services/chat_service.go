// api/services/chat_service.go

package services

import (
	"context"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService interface {
	SaveMessage(ctx context.Context, mID, uID primitive.ObjectID, content, name, profile string) (*models.ChatMessage, error)
	GetChatHistory(ctx context.Context, meetingID primitive.ObjectID, limit int64) ([]models.ChatMessage, error)
}

type chatService struct {
	chatRepo repositories.ChatRepository
	userRepo repositories.UserRepository
}

func NewChatService(cr repositories.ChatRepository, ur repositories.UserRepository) ChatService {
	return &chatService{chatRepo: cr, userRepo: ur}
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

func (s *chatService) GetChatHistory(ctx context.Context, meetingID primitive.ObjectID, limit int64) ([]models.ChatMessage, error) {
	return s.chatRepo.GetChatHistory(ctx, meetingID, limit)
}
