// api/services/auth_service.go

package services

import (
	"context"
	"errors"
	"time"

	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/models"
	"github.com/seojoonrp/bbiyong-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (string, error)
	Login(ctx context.Context, req models.LoginRequest) (string, error)
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
	CompleteProfile(ctx context.Context, userID string, req models.SetProfileRequest) error
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (string, error) {
	exists, _ := s.repo.FindByUsername(ctx, req.Username)
	if exists != nil {
		return "", errors.New("username already taken")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	user := models.User{
		Username:     req.Username,
		Password:     string(hashedPassword),
		Provider:     "local",
		IsProfileSet: false,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID.Hex())
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("incorrect password")
	}

	return utils.GenerateToken(user.ID.Hex())
}

func (s *authService) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	return user == nil, nil
}

func (s *authService) CompleteProfile(ctx context.Context, userID string, req models.SetProfileRequest) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id format")
	}

	// TODO : IsProfileSet이 false일 때만 업데이트되도록 수정

	// TODO : 다른 필드들 추가하기
	updates := bson.M{
		"nickname":       req.Nickname,
		"is_profile_set": true,
	}

	return s.repo.UpdateProfile(ctx, objID, updates)
}
