package product

import (
	"backend/server/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrMissingFile     error = errors.New("missing file")
	ErrMissingProduct  error
	ErrMissingRoomChat error = fmt.Errorf("you dont have any room chat,try created one")
	ErrMissingAccess   error = fmt.Errorf("you dont have any access to this chat")
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

func (s *service) GetMyProduct(ctx context.Context, userID int) ([]*GetProductResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	prc := []*GetProductResponse{}

	pr, err := s.Repository.GetProductByUserID(c, userID)
	if err != nil {
		return []*GetProductResponse{}, err
	}

	if len(pr) <= 0 {
		ErrMissingProduct = errors.New("you dont have any product yet")
		return []*GetProductResponse{}, ErrMissingProduct
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

func (s *service) CreateRoomChat(ctx context.Context, URC *CreateRoomChatRequest) (*CreateRoomChatResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	productID := strconv.Itoa(URC.ProductID)
	userID := strconv.Itoa(URC.UserID)

	rID, err := utils.GenerateRoomID(productID, userID, 21)
	if err != nil {
		return &CreateRoomChatResponse{}, err
	}

	id, roomID, err := s.Repository.CreateRoomChat(c, &UlasanRoomChatProduct{
		RoomID:    rID,
		UserID:    URC.UserID,
		ProductID: URC.ProductID,
		Username:  URC.Username,
	})
	if err != nil {
		return &CreateRoomChatResponse{}, err
	}

	return &CreateRoomChatResponse{
		ID:     id,
		RoomID: roomID,
	}, nil
}

func (s *service) CreateRoomChatMessage(ctx context.Context, URC *CreateRoomChatProductMessageReq) (*CreateRoomChatProductMessageRes, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check if the room exist return error if it isnt exist
	roomChat, err := s.Repository.GetRoomByProductID(c, URC.ProductID)
	if err != nil {
		return &CreateRoomChatProductMessageRes{}, err
	}

	if len(roomChat) < 1 {
		return &CreateRoomChatProductMessageRes{}, ErrMissingRoomChat
	}

	product, err := s.Repository.GetProductByID(c, URC.ProductID)
	if err != nil {
		return &CreateRoomChatProductMessageRes{}, err
	}

	// store into map for easier access
	roomMap := make(map[int]UlasanRoomChatProduct)
	for _, v := range roomChat {
		roomMap[v.ID] = *v
	}

	// check if roomID is valid
	message, err := validRoomChatID(roomMap, URC, *product)
	if err != nil {
		return &CreateRoomChatProductMessageRes{}, err
	}

	roomID, err := s.Repository.CreateRoomChatMessage(c, &UlasanRoomChatProductMessage{
		RoomID:    URC.RoomID,
		Message:   message,
		SendAt:    time.Now(),
		IsDeleted: false,
	})
	if err != nil {
		return &CreateRoomChatProductMessageRes{}, err
	}

	return &CreateRoomChatProductMessageRes{
		RoomID:  roomID,
		Message: message,
	}, nil
}

func validRoomChatID(rMap map[int]UlasanRoomChatProduct, urc *CreateRoomChatProductMessageReq, product Product) (string, error) {
	room, ok := rMap[urc.RoomID]
	if !ok {
		return "", fmt.Errorf("cannot find room id,have you ever created one ?")
	}

	if urc.Message == "" {
		urc.Message = "empty-string"
	}

	if urc.UserID != room.UserID && urc.UserID != product.UserID {
		return "", ErrMissingAccess
	}

	return urc.Message, nil
}
