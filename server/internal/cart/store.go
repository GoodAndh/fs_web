package cart

import (
	"context"
	"time"
)

type Cart struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProductID   int       `json:"product_id"`
	Status      string    `json:"status"`
	Total       int       `json:"total"`
	Price       float64   `json:"price"`
	ProductName string    `json:"product_name"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type CreateCartRequest struct {
	UserID    int    `json:"user_id" validate:"required"`
	ProductID int    `json:"product_id" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=paid wait"`
	Total     int    `json:"total" validate:"required"`
}

type CreateCartResponse struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	Status      string  `json:"status"`
	Total       int     `json:"total"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
	ProductName string  `json:"product_name"`
}

type UpdateCartRequest struct {
	ID        int    `json:"id" validate:"required"`
	UserID    int    `json:"user_id" validate:"required"`
	ProductID int    `json:"product_id" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=wait paid"`
	Total     int    `json:"total" validate:"required"`
}

type UpdateCartResponse struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	Status      string  `json:"status"`
	Total       int     `json:"total"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
	ProductName string  `json:"product_name"`
}

type GetCartResponse struct {
	UserID      int     `json:"user_id"`
	ProductID   int     `json:"product_id"`
	Status      string  `json:"status"`
	Total       int     `json:"total"`
	Price       float64 `json:"price"`
	ProductName string  `json:"product_name"`
}

type Repository interface {
	AddNewCart(ctx context.Context, cart *Cart) (*Cart, error)
	GetCartByProductID(ctx context.Context, userID, productID int) error
	UpdateCart(ctx context.Context, cart *Cart) error
	GetCartByID(ctx context.Context, ID int) (*Cart, error)
	GetCartByUserID(ctx context.Context, userID int) ([]*Cart, error)
}

type Service interface {
	AddNewCart(ctx context.Context, req *CreateCartRequest) (*CreateCartResponse, error)
	UpdateCart(ctx context.Context, req *UpdateCartRequest) (*UpdateCartResponse, error)
	GetCart(ctx context.Context,userID int)([]*GetCartResponse,error)
}
