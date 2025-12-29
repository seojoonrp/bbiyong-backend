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
	SaveSystemMessage(ctx context.Context, mID, uID primitive.ObjectID, eventType string) (*models.ChatMessage, error)
	CheckParticipation(ctx context.Context, mID, uID primitive.ObjectID) (bool, error)
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

func (s *chatService) SaveSystemMessage(ctx context.Context, mID, uID primitive.ObjectID, eventType string) (*models.ChatMessage, error) {
	user, err := s.userRepo.FindByID(ctx, uID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
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
		return nil, err
	}

	return msg, nil
}

func (s *chatService) verifyParticipant(ctx context.Context, mID, uID primitive.ObjectID) error {
	meeting, err := s.meetingRepo.FindByID(ctx, mID)
	if err != nil {
		return err
	}
	if meeting == nil {
		return errors.New("meeting not found")
	}

	for _, pID := range meeting.Participants {
		if pID == uID {
			return nil
		}
	}

	return errors.New("you are not a participant of the meeting")
}

func (s *chatService) CheckParticipation(ctx context.Context, mID, uID primitive.ObjectID) (bool, error) {
	err := s.verifyParticipant(ctx, mID, uID)
	if err != nil {
		if err.Error() == "you are not a participant of the meeting" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *chatService) GetChatHistory(ctx context.Context, mID, uID primitive.ObjectID, limit int64) ([]models.ChatMessage, error) {
	if err := s.verifyParticipant(ctx, mID, uID); err != nil {
		return nil, err
	}

	return s.chatRepo.GetChatHistory(ctx, mID, limit)
}
