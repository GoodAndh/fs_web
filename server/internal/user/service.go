package user

import (
	"backend/config"
	"backend/server/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUsernameAlrInUsed error = errors.New("username already in used")
	ErrEmailAlrInUsed    error = errors.New("email already in used")
	ErrMissingFile       error = errors.New("missing file")
)

type service struct {
	Repository
	timeout time.Duration
	*config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		Repository: repo,
		timeout:    time.Duration(2 * time.Second),
		Config:     cfg,
	}
}

func (s *service) CreateUsers(ctx context.Context, user *RegisUserRequest) (*RegisUserResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.Repository.GetUserByUsername(c, user.Username); err == nil {
		return &RegisUserResponse{}, ErrUsernameAlrInUsed
	}

	if _, err := s.Repository.GetUserByEmail(c, user.Email); err == nil {
		return &RegisUserResponse{}, ErrEmailAlrInUsed
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return &RegisUserResponse{}, err
	}

	userID, err := s.Repository.CreateUsers(c, &User{
		Username:    user.Username,
		Email:       user.Email,
		Password:    hashed,
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
	})
	if err != nil {
		return &RegisUserResponse{}, err
	}

	return &RegisUserResponse{
		ID:       strconv.Itoa(userID),
		Username: user.Username,
		Email:    user.Email,
	}, nil

}

type JWTclaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	cfg      string
	jwt.RegisteredClaims
}

func (s *service) LoginUsers(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByUsername(c, req.Username)
	if err != nil {
		if err == utils.ErrNotFound {
			return &LoginUserResponse{}, fmt.Errorf("invalid username or password")
		}
		return &LoginUserResponse{}, err
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return &LoginUserResponse{}, fmt.Errorf("invalid username or password")
	}

	jwtStx := JWTclaims{
		ID:               strconv.Itoa(user.ID),
		Username:         user.Username,
		cfg:              s.JWTSecretKey,
		RegisteredClaims: jwt.RegisteredClaims{Issuer: strconv.Itoa(user.ID), ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtStx)

	tokenString, err := token.SignedString([]byte(jwtStx.cfg))
	if err != nil {
		return &LoginUserResponse{}, err
	}

	return &LoginUserResponse{
		accessToken: tokenString,
		ID:          strconv.Itoa(user.ID),
		Username:    user.Username,
	}, nil

}

func (s *service) GetUserProfile(ctx context.Context, userID int) (*UserProfileResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	userProfile, err := s.Repository.GetUserProfile(c, userID)
	if err != nil {
		if err == utils.ErrNotFound {
			return &UserProfileResponse{}, fmt.Errorf("you dont have any profile")
		}
		return &UserProfileResponse{}, err
	}

	return &UserProfileResponse{
		ID:       strconv.Itoa(userProfile.ID),
		Url:      userProfile.Url,
		Captions: userProfile.Captions,
	}, nil
}

func (s *service) CreateUserProfile(ctx context.Context, req *UserProfileRequest) (*UserProfileResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.Repository.GetUserProfile(c, req.UserID); err == nil {
		return &UserProfileResponse{}, fmt.Errorf("you already have profile .do you mean update?")
	}

	id, err := s.Repository.CreateUserProfile(c, &UserProfile{
		UserID:   req.UserID,
		Url:      req.Url,
		Captions: req.Captions,
	})
	if err != nil {
		return &UserProfileResponse{}, err
	}

	return &UserProfileResponse{
		ID:       strconv.Itoa(id),
		Url:      req.Url,
		Captions: req.Captions,
	}, nil
}

func (s *service) UpdateUserProfile(ctx context.Context, req *UserProfileRequest) (*UpdateUserProfileResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.Repository.GetUserProfile(c, req.UserID); err != nil {
		if err == utils.ErrNotFound {
			return &UpdateUserProfileResponse{}, fmt.Errorf("you dont have any image yet,try create one")
		}
		return &UpdateUserProfileResponse{}, err
	}

	if err := s.Repository.UpdateUserProfile(c, &UserProfile{
		UserID:   req.UserID,
		Url:      req.Url,
		Captions: req.Captions,
	}); err != nil {
		return &UpdateUserProfileResponse{}, err
	}

	return &UpdateUserProfileResponse{
		Url:      req.Url,
		Captions: req.Captions,
	}, nil
}

func (s *service) GetUserByID(ctx context.Context, userID int) (*UserResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByID(c, userID)
	if err != nil {
		if err == utils.ErrNotFound {
			return &UserResponse{}, fmt.Errorf("id not found")
		}
		return &UserResponse{}, err
	}

	return &UserResponse{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
