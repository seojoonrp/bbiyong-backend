// api/services/friend_service.go

package services

import (
	"context"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendService interface {
	RequestFriend(ctx context.Context, userID, targetID string) error
	AcceptFriend(ctx context.Context, userID, friendshipID string) error
	ListFriends(ctx context.Context, userID string, status string) ([]models.FriendInfo, error)
}

type friendService struct {
	friendRepo repositories.FriendRepository
}

func NewFriendService(fr repositories.FriendRepository) FriendService {
	return &friendService{friendRepo: fr}
}

func (s *friendService) RequestFriend(ctx context.Context, userID, targetID string) error {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperr.InternalServerError("invalid user ID in token", err)
	}

	tID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return apperr.BadRequest("invalid target user ID format", err)
	}

	if uID == tID {
		return apperr.BadRequest("cannot friend yourself", nil)
	}

	existing, err := s.friendRepo.FindByUserIDs(ctx, uID, tID)
	if err != nil {
		return apperr.InternalServerError("failed to check existing friendship", err)
	}
	if existing != nil {
		return apperr.BadRequest("friend request already exists", nil)
	}

	friendship := &models.Friendship{
		ID:          primitive.NewObjectID(),
		RequesterID: uID,
		AddresseeID: tID,
		Status:      models.FriendStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.friendRepo.SendRequest(ctx, friendship)
	if err != nil {
		return apperr.InternalServerError("failed to send friend request", err)
	}

	return nil
}

func (s *friendService) AcceptFriend(ctx context.Context, userID, friendshipID string) error {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperr.InternalServerError("invalid user ID in token", err)
	}

	fID, err := primitive.ObjectIDFromHex(friendshipID)
	if err != nil {
		return apperr.BadRequest("invalid friendship ID format", err)
	}

	friendship, err := s.friendRepo.FindByID(ctx, fID)
	if err != nil {
		return apperr.InternalServerError("failed to fetch friendship", err)
	}
	if friendship == nil {
		return apperr.NotFound("friendship not found", nil)
	}
	if friendship.AddresseeID != uID {
		return apperr.Forbidden("you are not the addressee of the friend request", nil)
	}
	if friendship.Status != models.FriendStatusPending {
		return apperr.BadRequest("friendship is not in a pending state", nil)
	}

	err = s.friendRepo.UpdateStatus(ctx, fID, models.FriendStatusAccepted)
	if err != nil {
		return apperr.InternalServerError("failed to update friendship status", err)
	}

	return nil
}

func (s *friendService) ListFriends(ctx context.Context, userID string, status string) ([]models.FriendInfo, error) {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperr.InternalServerError("invalid user ID in token", err)
	}

	friendInfos, err := s.friendRepo.GetFriendList(ctx, uID, status)
	if err != nil {
		return nil, apperr.InternalServerError("failed to get friend list", err)
	}

	return friendInfos, nil
}
