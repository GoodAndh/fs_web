package orders

import (
	"backend/server/utils"
	"fmt"
	"strconv"

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
	router.Post("/order/create", h.MiddlewareWithJWT, h.createOrders)
	router.Get("/order/", h.MiddlewareWithJWT, h.getOrders)
}

func (h *Handler) createOrders(c *fiber.Ctx) error {

	var payload CreateOrderRequest
	if err := c.BodyParser(&payload); err != nil {
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
	payload.UserID = userID

	if errs := h.XValidator.Validate(&payload); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)
		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf("[%s:%v] need to implement '%s'", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "failed field on :", errMsg)
	}

	response, err := h.Service.CreateOrders(c.Context(), &payload)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success create orders", response)
}

func (h *Handler) getOrders(c *fiber.Ctx) error {
	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}
	response, err := h.Service.GetOrders(c.Context(), userID)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "order history:", response)
}
