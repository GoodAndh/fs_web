package orders

import (
	"backend/server/internal/product"
	"backend/server/utils"
	"context"
	"fmt"
	"time"
)

type service struct {
	repo        Repository
	repoProduct product.Repository
	timeout     time.Duration
}

func NewService(repo Repository, product product.Repository) Service {
	return &service{repo, product, time.Duration(2 * time.Second)}
}

func (s *service) GetOrders(ctx context.Context, userID int) ([]*GetOrdersResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	getO := []*GetOrdersResponse{}

	status, err := s.repo.GetStatusByUserID(c, userID)
	if err != nil {
		return []*GetOrdersResponse{}, err
	}

	if len(status) <= 0 {
		return []*GetOrdersResponse{}, fmt.Errorf("you dont have any order yet")
	}

	items, err := s.repo.GetItemsByUserID(c, userID)
	if err != nil {
		return []*GetOrdersResponse{}, err
	}

	if len(items) <= 0 {
		return []*GetOrdersResponse{}, fmt.Errorf("you dont have any order yet")

	}

	for i, stats := range status {
		var totalItems []int
		var idOrder []int
		for _, v := range items {
			totalItems = append(totalItems, v.Total)
			idOrder = append(idOrder, v.ID)
		}
		pr, err := s.repoProduct.GetProductByID(c, stats.ProductID)
		if err != nil {
			if err == utils.ErrNotFound {
				return []*GetOrdersResponse{}, fmt.Errorf("productID not found,something went wrong,actual response:%v ", err)
			}
			return []*GetOrdersResponse{}, err
		}
		ord := &GetOrdersResponse{
			Id:          idOrder[i],
			OrderID:     stats.ID,
			UserID:      userID,
			ProductID:   stats.ProductID,
			ProductName: pr.Name,
			Total:       totalItems[i],
			TotalPrice:  stats.TotalPrice,
			Status:      stats.Status,
		}
		getO = append(getO, ord)
	}

	if len(getO) <= 0 {
		return []*GetOrdersResponse{}, fmt.Errorf("you dont have any order yet")
	}

	return getO, nil
}

func (s *service) CreateOrders(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check the product if exist
	pr, err := s.repoProduct.GetProductByID(c, req.ProductID)
	if err != nil {
		if err == utils.ErrNotFound {
			return &CreateOrderResponse{}, fmt.Errorf("product id :%d not found ", req.ProductID)
		}
		return &CreateOrderResponse{}, err
	}

	// check if the product alr in order
	oStatus, err := s.repo.GetOrderStatus(c, req.UserID, req.ProductID)
	if err == nil && oStatus.Status == "wait" {
		return &CreateOrderResponse{}, fmt.Errorf("you alr have the same order ,do you mean change?")
	}

	// check if in stock or is the user own it self
	if err := checkinStock(pr, req.UserID, req.Total); err != nil {
		return &CreateOrderResponse{}, err
	}

	idStatus, idItems, err := s.repo.CreateOrders(c, &OrderStatus{
		ProductID:  req.ProductID,
		UserID:     req.UserID,
		Status:     "wait",
		TotalPrice: pr.Price * float64(req.Total),
	}, &OrderItems{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Total:     req.Total,
		Price:     pr.Price,
	})
	if err != nil {
		return &CreateOrderResponse{}, err
	}

	return &CreateOrderResponse{Id: idItems, OrderID: idStatus, UserID: req.UserID, ProductID: req.UserID, Total: req.Total, TotalPrice: pr.Price * float64(req.Total)}, nil

}

func checkinStock(pr *product.Product, userID, total int) error {

	if pr.UserID == userID {
		return fmt.Errorf("you cant order your own product")
	}

	if total > pr.Stock {
		return fmt.Errorf("order limit,remaining stock '%d'", pr.Stock)
	}

	return nil
}
