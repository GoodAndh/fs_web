package user

import (
	"backend/config"
	"backend/server/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		return &RegisUserResponse{}, fmt.Errorf("username already in used")
	}

	if _, err := s.Repository.GetUserByEmail(c, user.Email); err == nil {
		return &RegisUserResponse{}, fmt.Errorf("email already in used")
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
