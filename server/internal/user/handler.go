package user

import (
	"backend/server/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service
	utils.XValidator
}

func NewHandler(service Service, validate utils.XValidator) *Handler {
	return &Handler{service, validate}
}

func (h *Handler) RegisterRoute(router fiber.Router) {
	router.Post("/signup", h.createUsers)
	router.Post("/signin", h.loginUsers)
	router.Get("/signout", h.logoutUsers)
}

func (h *Handler) createUsers(c *fiber.Ctx) error {
	var Req RegisUserRequest
	if err := c.BodyParser(&Req); err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	if errs := h.XValidator.Validate(&Req); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)
		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf("[%s:%v] need to implement '%s'", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "failed field on :", errMsg)
	}

	response, err := h.Service.CreateUsers(c.Context(), &Req)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success create users", response)

}

func (h *Handler) loginUsers(c *fiber.Ctx) error {
	var Req LoginUserRequest
	if err := c.BodyParser(&Req); err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	response, err := h.LoginUsers(c.Context(), &Req)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    response.accessToken,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
		MaxAge:   3600,
		Path:     "/",
	})

	return utils.WriteJson(c, 200, "login success", &LoginUserResponse{
		ID:       response.ID,
		Username: response.Username,
	})

}

func (h *Handler) logoutUsers(c *fiber.Ctx) error {
	if c.Cookies("jwt") != "" {
		c.Cookie(&fiber.Cookie{
			Name:    "jwt",
			Expires: time.Now().Add(-1),
		})
		return utils.WriteJson(c, 200, "cookies deleted", nil)
	}
	return utils.WriteJson(c, 401, "unauthorized", nil)
}
