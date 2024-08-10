package product

import (
	"context"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Url       string `json:"url"`
	Captions  string `json:"captions"`
}

type CreateProductRequest struct {
	UserID      int     `json:"user_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Stock       int     `json:"stock" validate:"required"`
}

type CreateProductResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type GetProductResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type CreateProductImageRequest struct {
	ProductID int    `json:"product_id" validate:"required"`
	Url       string `json:"url" validate:"required"`
	Captions  string `json:"captions" validate:"required"`
}

type ProductImageResponse struct {
	Url      string `json:"url"`
	Captions string `json:"captions"`
}

type Repository interface {
	CreateProduct(ctx context.Context, product *Product) (int, error)
	GetAllProduct(ctx context.Context) ([]*Product, error)
	GetProductByID(ctx context.Context, id int) (*Product, error)
	GetProductByName(ctx context.Context, name string) (*Product, error)
	GetProductByUserID(ctx context.Context, userID int) ([]*Product, error)
	CreateProductImage(ctx context.Context, img *ProductImage) error
	GetProductImage(ctx context.Context, productID int) ([]*ProductImage, error)
}

type Service interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductResponse, error)
	GetAllProduct(ctx context.Context) ([]*GetProductResponse, error)
	GetProductByID(ctx context.Context, id int) (*GetProductResponse, error)
	GetProductByName(ctx context.Context, name string) (*GetProductResponse, error)
	CreateProductImage(ctx context.Context,userID int, req *CreateProductImageRequest) (*ProductImageResponse, error)
	GetProductImage(ctx context.Context, productID int)([]*ProductImageResponse,error)
	GetMyProduct(ctx context.Context,userID int)([]*GetProductResponse,error)
}
