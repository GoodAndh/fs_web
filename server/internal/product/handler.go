package product

import (
	"backend/server/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service
	utils.XValidator
	*utils.MiddlewareStruct
}

func NewHandler(s Service, x utils.XValidator, m *utils.MiddlewareStruct) *Handler {
	return &Handler{s, x, m}
}

func (h *Handler) RegisterRoute(router fiber.Router) {
	router.Post("/product/create", h.MiddlewareWithJWT, h.createProduct)
	router.Get("/product", h.getAllProduct)
	router.Get("/product/id/:id", h.getProductByID)
	router.Get("/product/search/:name", h.getProductByName)
}

func (h *Handler) createProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}
	req.UserID = userID

	if errs := h.XValidator.Validate(&req); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)
		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf("[%s:%v] need to implement '%s'", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "failed field on :", errMsg)
	}

	response, err := h.Service.CreateProduct(c.Context(), &req)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success create product", response)

}

func (h *Handler) getAllProduct(c *fiber.Ctx) error {

	product, err := h.Service.GetAllProduct(c.Context())
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "all product", product)

}

func (h *Handler) getProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}
	product, err := h.Service.GetProductByID(c.Context(), id)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.WriteJson(c, 400, fmt.Sprintf("[id:%v] not found ,make sure the id was int or not equal to string", id), nil)
		}
		return utils.WriteJson(c, 400, err.Error(), nil)
	}
	return utils.WriteJson(c, 200, "product by id", product)
}

func (h *Handler) getProductByName(c *fiber.Ctx) error {
	name := c.Params("name", "not found")
	newName:=strings.ReplaceAll(name,"%20"," ")
	product, err := h.Service.GetProductByName(c.Context(), newName)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.WriteJson(c, 400, fmt.Sprintf("[name:%v] not found ", newName), nil)
		}
		return utils.WriteJson(c, 400, err.Error(), nil)
	}
	return utils.WriteJson(c, 200, "product by id", product)
}
