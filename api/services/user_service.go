// api/services/user_service.go

package services

import (
	"context"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(ur repositories.UserRepository) UserService {
	return &userService{userRepo: ur}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	uID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperr.InternalServerError("invalid user ID in token", err)
	}

	user, err := s.userRepo.FindByID(ctx, uID)
	if err != nil {
		return nil, apperr.InternalServerError("failed to fetch user by id", err)
	}
	if user == nil {
		return nil, apperr.NotFound("user not found", nil)
	}
	return user, nil
}
