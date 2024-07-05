package cart

import (
	"backend/server/utils"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service
	utils.XValidator
	utils.MiddlewareStruct
}

func NewHandler(s Service, x utils.XValidator, m utils.MiddlewareStruct) *Handler {
	return &Handler{s, x, m}
}

func (h *Handler) RegisterRoute(router fiber.Router) {
	router.Post("/cart/create/", h.MiddlewareWithJWT, h.createCart)
	router.Post("/cart/update/", h.MiddlewareWithJWT, h.updateCart)
	router.Get("/cart/", h.MiddlewareWithJWT, h.getCart)
}

func (h *Handler) createCart(c *fiber.Ctx) error {

	var req CreateCartRequest
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
			if err.Value == "" {
				err.Value = "empty value"
			}
			errMsg = append(errMsg, fmt.Sprintf("[%s:%v] need to implement '%s' ", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "failed field on:", errMsg)
	}

	response, err := h.Service.AddNewCart(c.Context(), &req)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}
	return utils.WriteJson(c, 200, "success add new cart", response)
}

func (h *Handler) updateCart(c *fiber.Ctx) error {
	var req UpdateCartRequest
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

			errMsg = append(errMsg, fmt.Sprintf("[%s:%v] need to implement '%s' ", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "failed field on:", errMsg)
	}

	response, err := h.Service.UpdateCart(c.Context(), &req)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success update cart", response)

}

func (h *Handler) getCart(c *fiber.Ctx) error {

	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	response, err := h.Service.GetCart(c.Context(), userID)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, fmt.Sprintf("you have '%d' cart", len(response)), response)
}
