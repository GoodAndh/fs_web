package cart

import (
	"backend/server/internal/product"
	"backend/server/utils"
	"context"
	"fmt"
	"time"
)

type service struct {
	Repository
	pRepo   product.Repository
	timeout time.Duration
}

func NewService(repo Repository, pRepo product.Repository) Service {
	return &service{
		Repository: repo,
		pRepo:      pRepo,
		timeout:    time.Duration(2 * time.Second),
	}
}

func (s *service) GetCart(ctx context.Context, userID int) ([]*GetCartResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cart, err := s.Repository.GetCartByUserID(c, userID)
	if err != nil {
		return []*GetCartResponse{}, err
	}

	if len(cart) <= 0 {
		return []*GetCartResponse{}, fmt.Errorf("u dont have any cart")
	}

	cartSlice := make([]*GetCartResponse, 0, len(cart))

	for _, v := range cart {
		cart := &GetCartResponse{
			UserID:      v.UserID,
			ProductID:   v.ProductID,
			Status:      v.Status,
			Total:       v.Total,
			Price:       v.Price,
			ProductName: v.ProductName,
		}
		cartSlice = append(cartSlice, cart)
	}

	return cartSlice, nil

}

func (s *service) UpdateCart(ctx context.Context, req *UpdateCartRequest) (*UpdateCartResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// validate if the cart is exist

	if _, err := s.Repository.GetCartByID(c, req.ID); err != nil {
		if err == utils.ErrNotFound {
			return &UpdateCartResponse{}, fmt.Errorf("kau belum memiliki cart dengan id: '%d'", req.ID)
		}
		return &UpdateCartResponse{}, err
	}

	if err := s.Repository.GetCartByProductID(c, req.UserID, req.ProductID); err != nil {
		if err == utils.ErrNotFound {
			return &UpdateCartResponse{}, fmt.Errorf("kau belum memiliki cart dengan productID: '%d'", req.ProductID)
		}
		return &UpdateCartResponse{}, err
	}

	// check if the product ID is exist
	prd, err := s.pRepo.GetProductByID(c, req.ProductID)
	if err != nil {
		if err == utils.ErrNotFound {
			return &UpdateCartResponse{}, fmt.Errorf("productID '%d' ini belum ada ", req.ProductID)
		}
		return &UpdateCartResponse{}, err
	}

	// check the remaining stock from product
	if err := checkInStock(prd, req.Total); err != nil {
		return &UpdateCartResponse{}, err
	}

	if err := s.Repository.UpdateCart(c, &Cart{
		ID:          req.ID,
		UserID:      req.UserID,
		ProductID:   req.ProductID,
		Status:      req.Status,
		Total:       req.Total,
		Price:       prd.Price,
		ProductName: prd.Name,
		LastUpdated: time.Now(),
	}); err != nil {
		return &UpdateCartResponse{}, err
	}

	return &UpdateCartResponse{
		ID:          req.ID,
		ProductID:   req.ProductID,
		Status:      req.Status,
		Total:       req.Total,
		Price:       prd.Price,
		TotalPrice:  calculateTotalPrice(prd, req.Total),
		ProductName: prd.Name,
	}, nil
}

func (s *service) AddNewCart(ctx context.Context, req *CreateCartRequest) (*CreateCartResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if prd, err := s.pRepo.GetProductByUserID(c, req.UserID); err == nil && len(prd) > 0 {
		return &CreateCartResponse{}, fmt.Errorf("you cant add your own item to the cart")
	}

	prd, err := s.pRepo.GetProductByID(c, req.ProductID)
	if err != nil {
		return &CreateCartResponse{}, err
	}

	if err := s.Repository.GetCartByProductID(c, req.UserID, req.ProductID); err == nil {
		return &CreateCartResponse{}, fmt.Errorf("you already have the same cart. do you mean change?")
	}

	// check if product is found or in stock
	if err := checkInStock(prd, req.Total); err != nil {
		return &CreateCartResponse{}, err
	}

	cart, err := s.Repository.AddNewCart(c, &Cart{
		UserID:      req.UserID,
		ProductID:   req.ProductID,
		Status:      req.Status,
		Total:       req.Total,
		Price:       prd.Price,
		CreatedAt:   time.Now(),
		ProductName: prd.Name,
		LastUpdated: time.Now(),
	})
	if err != nil {
		return &CreateCartResponse{}, err
	}

	return &CreateCartResponse{
		ID:          cart.ID,
		ProductID:   cart.ProductID,
		Status:      cart.Status,
		Total:       cart.Total,
		Price:       cart.Price,
		TotalPrice:  calculateTotalPrice(prd, req.Total),
		ProductName: cart.ProductName,
	}, nil

}

func checkInStock(prd *product.Product, req int) error {

	if prd.Stock < 0 {
		return fmt.Errorf("remaining stock left : '%d'", prd.Stock)
	}

	if req > prd.Stock {
		return fmt.Errorf("your requested item its higher than remaining stock : '%d'", prd.Stock)
	}

	return nil
}

func calculateTotalPrice(prd *product.Product, req int) float64 {
	total := prd.Price * float64(req)
	return total
}
