// api/services/friend_service.go

package services

import (
	"context"
	"errors"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendService interface {
	RequestFriend(ctx context.Context, reqID, addID primitive.ObjectID) error
	AcceptFriend(ctx context.Context, fID, uID primitive.ObjectID) error
	ListFriends(ctx context.Context, userID primitive.ObjectID) ([]models.Friendship, error)
}

type friendService struct {
	friendRepo repositories.FriendRepository
}

func NewFriendService(fr repositories.FriendRepository) FriendService {
	return &friendService{friendRepo: fr}
}

func (s *friendService) RequestFriend(ctx context.Context, reqID, addID primitive.ObjectID) error {
	if reqID == addID {
		return errors.New("cannot friend yourself")
	}

	existing, _ := s.friendRepo.FindByUserIDs(ctx, reqID, addID)
	if existing != nil {
		return errors.New("friendship already exists or is pending")
	}

	friendship := &models.Friendship{
		ID:          primitive.NewObjectID(),
		RequesterID: reqID,
		AddresseeID: addID,
		Status:      models.FriendStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.friendRepo.SendRequest(ctx, friendship)
}

func (s *friendService) AcceptFriend(ctx context.Context, fID, uID primitive.ObjectID) error {
	friendship, err := s.friendRepo.FindByID(ctx, fID)
	if err != nil {
		return err
	}
	if friendship == nil {
		return errors.New("friendship not found")
	}
	if friendship.AddresseeID != uID {
		return errors.New("you are not the addressee of the friend request")
	}
	if friendship.Status != models.FriendStatusPending {
		return errors.New("friendship is not in a pending state")
	}

	return s.friendRepo.UpdateStatus(ctx, fID, models.FriendStatusAccepted)
}

func (s *friendService) ListFriends(ctx context.Context, uID primitive.ObjectID) ([]models.Friendship, error) {
	return s.friendRepo.GetFriends(ctx, uID)
}
