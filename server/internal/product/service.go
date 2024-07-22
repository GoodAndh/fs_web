package product

import (
	"backend/server/utils"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrMissingFile error = errors.New("missing file")
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		Repository: r,
		timeout:    time.Duration(2 * time.Second),
	}
}

func (s *service) CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	prID, err := s.Repository.CreateProduct(c, &Product{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
	})
	if err != nil {
		return &CreateProductResponse{}, err
	}

	return &CreateProductResponse{
		ID:          prID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}, nil
}

func (s *service) GetAllProduct(ctx context.Context) ([]*GetProductResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	prc := []*GetProductResponse{}

	pr, err := s.Repository.GetAllProduct(c)
	if err != nil {
		return []*GetProductResponse{}, err
	}
	for _, v := range pr {
		product := &GetProductResponse{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Stock:       v.Stock,
		}
		prc = append(prc, product)
	}
	return prc, nil
}

func (s *service) GetProductByID(ctx context.Context, id int) (*GetProductResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	product, err := s.Repository.GetProductByID(c, id)
	if err != nil {
		return &GetProductResponse{}, err
	}
	return &GetProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil

}

func (s *service) GetProductByName(ctx context.Context, name string) (*GetProductResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	product, err := s.Repository.GetProductByName(c, name)
	if err != nil {
		return &GetProductResponse{}, err
	}
	return &GetProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *service) CreateProductImage(ctx context.Context, userID int, req *CreateProductImageRequest) (*ProductImageResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if pr, err := s.Repository.GetProductByID(c, req.ProductID); err != nil {
		if err == utils.ErrNotFound {
			return &ProductImageResponse{}, fmt.Errorf("product id : '%d' not found", req.ProductID)
		}
		return &ProductImageResponse{}, err
	} else if pr.UserID != userID {
		return &ProductImageResponse{}, fmt.Errorf("its not your product")
	}

	if err := s.Repository.CreateProductImage(c, &ProductImage{
		ProductID: req.ProductID,
		Url:       req.Url,
		Captions:  req.Captions,
	}); err != nil {
		return &ProductImageResponse{}, err
	}

	return &ProductImageResponse{
		Url:      req.Url,
		Captions: req.Captions,
	}, nil
}

func (s *service) GetProductImage(ctx context.Context, productID int) ([]*ProductImageResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	pm := []*ProductImageResponse{}

	if _, err := s.Repository.GetProductByID(c, productID); err != nil {
		if err == utils.ErrNotFound {
			return []*ProductImageResponse{}, fmt.Errorf("product id : '%d' not found", productID)
		}
		return []*ProductImageResponse{}, err
	}

	image, err := s.Repository.GetProductImage(c, productID)
	if err != nil {
		return []*ProductImageResponse{}, err
	}

	if len(image) <= 0 {
		return []*ProductImageResponse{}, fmt.Errorf("you dont have any image of this product id '%d'", productID)
	}

	for _, v := range image {
		pr := ProductImageResponse{
			Url:      v.Url,
			Captions: v.Captions,
		}
		pm = append(pm, &pr)
	}

	return pm, nil
}
