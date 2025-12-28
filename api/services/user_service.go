// api/services/user_service.go

package services

import (
	"context"
	"errors"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(ur repositories.UserRepository) UserService {
	return &userService{userRepo: ur}
}

func (s *userService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
