package orders

import (
	"context"
	"database/sql"
)

type OrderStatus struct {
	ID         int     `json:"id"`
	ProductID  int     `json:"product_id"`
	UserID     int     `json:"user_id"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
}

type OrderItems struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Total     int     `json:"total"`
	Price     float64 `json:"price"`
}

type CreateOrderRequest struct {
	UserID    int `json:"user_id" validate:"required"`
	ProductID int `json:"product_id" validate:"required"`
	Total     int `json:"total" validate:"required"`
}

type CreateOrderResponse struct {
	Id         int     `json:"id"`
	OrderID    int     `json:"order_id"`
	UserID     int     `json:"user_id"`
	ProductID  int     `json:"product_id"`
	Total      int     `json:"total"`
	TotalPrice float64 `json:"total_price"`
}

type GetOrdersResponse struct {
	Id          int     `json:"id"`
	OrderID     int     `json:"order_id"`
	UserID      int     `json:"user_id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Total       int     `json:"total"`
	TotalPrice  float64 `json:"total_price"`
	Status      string  `json:"status"`
}

type Repository interface {
	GetOrderStatus(ctx context.Context, userID, ProductID int) (*OrderStatus, error)
	CreateOrders(ctx context.Context, ord *OrderStatus, ordItems *OrderItems) (int, int, error)
	GetOrderItems(ctx context.Context, userID, productID int) (*OrderItems, error)
	GetStatusByUserID(ctx context.Context, userID int) ([]*OrderStatus, error)
	GetItemsByUserID(ctx context.Context, userID int) ([]*OrderItems, error)
}

type Service interface {
	CreateOrders(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrders(ctx context.Context, userID int)([]*GetOrdersResponse,error)
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
