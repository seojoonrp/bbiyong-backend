// api/services/auth_service.go

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/apperr"
	"github.com/seojoonrp/bbiyong-backend/config"
	"github.com/seojoonrp/bbiyong-backend/models"
	"github.com/seojoonrp/bbiyong-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, req models.LoginRequest) (string, *models.User, error)
	LoginWithGoogle(ctx context.Context, idToken string) (bool, string, *models.User, error)
	LoginWithKakao(ctx context.Context, accessToken string) (bool, string, *models.User, error)
	LoginWithApple(ctx context.Context, identityToken string) (bool, string, *models.User, error)
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
	CompleteProfile(ctx context.Context, userID string, req models.SetProfileRequest) error
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) error {
	if len(req.Username) < 3 || len(req.Username) > 15 {
		return apperr.BadRequest("username must be between 3 and 15 characters", nil)
	}

	exists, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return apperr.InternalServerError("failed to fetch user by username", err)
	}
	if exists != nil {
		return apperr.BadRequest("username already exists", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return apperr.InternalServerError("failed to hash password", err)
	}

	user := models.User{
		Username:     req.Username,
		Password:     string(hashedPassword),
		Nickname:     "",
		ProfileURI:   "",
		Age:          -1,
		Gender:       "",
		Level:        1,
		Location:     models.Location{},
		RegionName:   "",
		Provider:     models.ProviderLocal,
		IsProfileSet: false,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return apperr.InternalServerError("failed to create user", err)
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (string, *models.User, error) {
	user, usernameErr := s.userRepo.FindByUsername(ctx, req.Username)
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if user == nil || usernameErr != nil || passwordErr != nil {
		return "", nil, apperr.Unauthorized("invalid username or password", nil)
	}

	signedToken, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", nil, apperr.InternalServerError("failed to generate token", err)
	}

	return signedToken, user, nil
}

func (s *authService) loginWithSocial(ctx context.Context, provider string, socialID string, email string) (bool, string, *models.User, error) {
	targetUsername := utils.GenerateHashUsername(provider, socialID)
	isNew := false

	user, err := s.userRepo.FindByUsername(ctx, targetUsername)
	if err != nil {
		return false, "", nil, apperr.InternalServerError("failed to fetch user by username", err)
	}

	if user == nil {
		isNew = true
		user = &models.User{
			Username:     targetUsername,
			Nickname:     "",
			ProfileURI:   "",
			Age:          -1,
			Gender:       "",
			Level:        1,
			Location:     models.Location{},
			RegionName:   "",
			Provider:     provider,
			SocialID:     socialID,
			IsProfileSet: false,
			CreatedAt:    time.Now(),
		}
		if email != "" {
			user.SocialEmail = email
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return false, "", nil, apperr.InternalServerError("failed to create user", err)
		}
	}

	signedToken, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return false, "", nil, apperr.InternalServerError("failed to generate token", err)
	}

	return isNew, signedToken, user, nil
}

func (s *authService) LoginWithGoogle(ctx context.Context, idToken string) (bool, string, *models.User, error) {
	webClientID := config.AppConfig.GoogleWebClientID

	payload, err := idtoken.Validate(context.Background(), idToken, webClientID)
	if err != nil {
		return false, "", nil, apperr.Unauthorized("invalid Google ID token", err)
	}

	socialID := payload.Subject
	email, _ := payload.Claims["email"].(string)

	return s.loginWithSocial(ctx, models.ProviderGoogle, socialID, email)
}

func (s *authService) LoginWithKakao(ctx context.Context, accessToken string) (bool, string, *models.User, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return false, "", nil, apperr.ServiceUnavailable("kakao api server unreachable", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return false, "", nil, apperr.Unauthorized("expired or invalid kakao token", nil)
	} else if resp.StatusCode != http.StatusOK {
		return false, "", nil, apperr.InternalServerError("kakao api returned error status", fmt.Errorf("status: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	var kakaoRes struct {
		ID           int64 `json:"id"`
		KakaoAccount struct {
			Email   string `json:"email"`
			Profile struct {
				Nickname string `json:"nickname"`
			} `json:"profile"`
		} `json:"kakao_account"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&kakaoRes); err != nil {
		return false, "", nil, apperr.InternalServerError("failed to decode Kakao user info", err)
	}

	socialID := strconv.FormatInt(kakaoRes.ID, 10)
	email := kakaoRes.KakaoAccount.Email

	return s.loginWithSocial(ctx, models.ProviderKakao, socialID, email)
}

func (s *authService) verifyAppleToken(identityToken string, clientID string) (jwt.MapClaims, error) {
	appleJWKSURL := "https://appleid.apple.com/auth/keys"

	k, err := keyfunc.NewDefault([]string{appleJWKSURL})
	if err != nil {
		return nil, apperr.InternalServerError("failed to create keyfunc", err)
	}

	token, err := jwt.Parse(identityToken, k.Keyfunc)
	if err != nil {
		return nil, apperr.InternalServerError("invalid token", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["iss"] != "https://appleid.apple.com" {
			return nil, apperr.Unauthorized("invalid issuer", nil)
		}
		if claims["aud"] != clientID {
			return nil, apperr.Unauthorized("invalid audience", nil)
		}
		return claims, nil
	}

	return nil, apperr.Unauthorized("invalid token claims", nil)
}

func (s *authService) LoginWithApple(ctx context.Context, identityToken string) (bool, string, *models.User, error) {
	clientID := config.AppConfig.AppleBundleID
	claims, err := s.verifyAppleToken(identityToken, clientID)
	if err != nil {
		return false, "", nil, err
	}

	socialID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)

	return s.loginWithSocial(ctx, models.ProviderApple, socialID, email)
}

func (s *authService) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return false, apperr.InternalServerError("failed to fetch user by username", err)
	}
	return user == nil, nil
}

func (s *authService) CompleteProfile(ctx context.Context, userID string, req models.SetProfileRequest) error {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperr.InternalServerError("invalid user ID in token", err)
	}

	updates := bson.M{
		"nickname":       req.Nickname,
		"profile_uri":    req.ProfileURI,
		"age":            req.Age,
		"gender":         req.Gender,
		"location":       req.Location,
		"region_name":    req.RegionName,
		"is_profile_set": true,
	}

	success, err := s.userRepo.CompleteProfile(ctx, uID, updates)
	if err != nil {
		return apperr.InternalServerError("failed to complete profile", err)
	}
	if !success {
		return apperr.BadRequest("profile is already set", nil)
	}

	return nil
}
