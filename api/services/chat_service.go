// api/services/chat_service.go

package services

import (
	"context"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService interface {
	SaveMessage(ctx context.Context, meetingID, userID string, content, name, profile string) (*models.ChatMessage, error)
	SaveSystemMessage(ctx context.Context, meetingID, userID string, eventType string) (*models.ChatMessage, error)
	GetChatHistory(ctx context.Context, meetingID string, limit int64) ([]models.ChatMessage, error)
}

type chatService struct {
	chatRepo    repositories.ChatRepository
	userRepo    repositories.UserRepository
	meetingRepo repositories.MeetingRepository
}

func NewChatService(cr repositories.ChatRepository, ur repositories.UserRepository, mr repositories.MeetingRepository) ChatService {
	return &chatService{chatRepo: cr, userRepo: ur, meetingRepo: mr}
}

func (s *chatService) SaveMessage(ctx context.Context, meetingID, userID string, content, name, profile string) (*models.ChatMessage, error) {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperr.InternalServerError("invalid user ID in token", err)
	}

	mID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		return nil, apperr.BadRequest("invalid meeting ID format", err)
	}

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

	err = s.chatRepo.SaveMessage(ctx, msg)
	if err != nil {
		return nil, apperr.InternalServerError("failed to save message", err)
	}

	return msg, nil
}

func (s *chatService) SaveSystemMessage(ctx context.Context, meetingID, userID string, eventType string) (*models.ChatMessage, error) {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperr.InternalServerError("invalid user ID in token", err)
	}

	mID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		return nil, apperr.BadRequest("invalid meeting ID format", err)
	}

	user, err := s.userRepo.FindByID(ctx, uID)
	if err != nil {
		return nil, apperr.InternalServerError("failed to fetch user by ID", err)
	}
	if user == nil {
		return nil, apperr.NotFound("user not found", nil)
	}

	var content, chatType string
	switch eventType {
	case models.EventJoinMeeting:
		content = user.Nickname + "님이 참여했습니다."
		chatType = models.ChatTypeJoin
	case models.EventLeaveMeeting:
		content = user.Nickname + "님이 나갔습니다."
		chatType = models.ChatTypeLeave
	default:
		content = "알 수 없는 이벤트가 발생했습니다."
		chatType = "unknown"
	}

	msg := &models.ChatMessage{
		ID:               primitive.NewObjectID(),
		MeetingID:        mID,
		SenderID:         uID,
		SenderName:       "System",
		SenderProfileURI: "",
		Content:          content,
		Type:             chatType,
		CreatedAt:        time.Now(),
	}

	err = s.chatRepo.SaveMessage(ctx, msg)
	if err != nil {
		return nil, apperr.InternalServerError("failed to save system message", err)
	}

	return msg, nil
}

func (s *chatService) GetChatHistory(ctx context.Context, meetingID string, limit int64) ([]models.ChatMessage, error) {
	mID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		return nil, apperr.BadRequest("invalid meeting ID format", err)
	}

	if limit <= 0 {
		return nil, apperr.BadRequest("limit must be greater than zero", nil)
	}
	if limit > 100 {
		return nil, apperr.BadRequest("cannot fetch more than 100 messages at once", nil)
	}

	history, err := s.chatRepo.GetChatHistory(ctx, mID, limit)
	if err != nil {
		return nil, apperr.InternalServerError("failed to get chat history", err)
	}

	return history, nil
}
