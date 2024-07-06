package user

import (
	"context"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type UserProfile struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Url      string `json:"url"`
	Captions string `json:"captions"`
}

type RegisUserRequest struct {
	Username  string `json:"username" validate:"required,min=8"`
	Email     string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required,min=8"`
	VPassword string `json:"vpassword" validate:"required,eqfield=Password"`
}

type RegisUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	accessToken string
	ID          string `json:"id"`
	Username    string `json:"username"`
}

type CreateUserProfileRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	Url      string `json:"url" validate:"required"`
	Captions string `json:"captions" validate:"required"`
}

type UserProfileResponse struct {
	ID       string `json:"id"`
	Url      string `json:"url"`
	Captions string `json:"captions"`
}

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUsers(ctx context.Context, user *User) (int, error)
	GetUserProfile(ctx context.Context, userID int) (*UserProfile, error)
	CreateUserProfile(ctx context.Context, user *UserProfile) (int, error)
}

type Service interface {
	CreateUsers(ctx context.Context, user *RegisUserRequest) (*RegisUserResponse, error)
	LoginUsers(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error)
	GetUserProfile(ctx context.Context, userID int) (*UserProfileResponse, error)
	CreateUserProfile(ctx context.Context, req *CreateUserProfileRequest) (*UserProfileResponse, error)
}
