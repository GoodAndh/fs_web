package user

import (
	"backend/server/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service
	utils.XValidator
	utils.MiddlewareStruct
}

func NewHandler(service Service, validate utils.XValidator, m utils.MiddlewareStruct) *Handler {
	return &Handler{service, validate, m}
}

func (h *Handler) RegisterRoute(router fiber.Router) {
	router.Post("/signup", h.createUsers)
	router.Post("/signin", h.loginUsers)
	router.Get("/signout", h.logoutUsers)
	router.Get("/user/profile", h.MiddlewareWithJWT, h.getProfile)
	router.Post("/user/profile", h.MiddlewareWithJWT, h.createProfile)
	router.Post("/user/profile/update", h.MiddlewareWithJWT, h.updateProfile)
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

func (h *Handler) getProfile(c *fiber.Ctx) error {
	idFromContext, ok := c.Locals("userID").(string)
	if !ok {
		return utils.WriteJson(c, 401, "unauthorized", nil)
	}

	userID, err := strconv.Atoi(idFromContext)
	if err != nil {
		return utils.WriteJson(c, 401, "unauthorized", nil)
	}

	response, err := h.Service.GetUserProfile(c.Context(), userID)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "profile image:", response)
}

func (h *Handler) createProfile(c *fiber.Ctx) error {
	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	url := c.Query("url")
	captions := c.Query("captions")
	payload := UserProfileRequest{
		UserID:   userID,
		Url:      url,
		Captions: captions,
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	files := form.File["file"]
	if len(files) == 0 {
		return utils.WriteJson(c, 400, "missing file", nil)
	}

	if len(files) > 1 {
		return utils.WriteJson(c, 400, fmt.Sprintf("only accept 1 file ,got '%d'", len(files)), nil)
	}

	if errs := h.XValidator.Validate(&payload); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)
		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf("[%s:%v] need to implement '%s'", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "failed field on :", errMsg)
	}

	file := files[0]

	fileExt := filepath.Ext(file.Filename)

	dir, err := os.Getwd()
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	newUrl := strings.ReplaceAll(payload.Url, " ", "_")

	file.Filename = newUrl + fileExt
	dest := filepath.Join(dir, "img", file.Filename)

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return utils.WriteJson(c, 500, "mkdir error:"+err.Error(), nil)
	}

	if err := c.SaveFile(file, dest); err != nil {
		return utils.WriteJson(c, 500, "savefile err:"+err.Error(), nil)
	}

	payload.Url = file.Filename

	response, err := h.Service.CreateUserProfile(c.Context(), &payload)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success add profile", response)
}

func (h *Handler) updateProfile(c *fiber.Ctx) error {
	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	url := c.Query("url")
	captions := c.Query("captions")
	payload := UserProfileRequest{
		UserID:   userID,
		Url:      url,
		Captions: captions,
	}
	form, err := c.MultipartForm()
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	files := form.File["file"]
	if len(files) == 0 {
		return utils.WriteJson(c, 400, "missing file", nil)
	}

	if len(files) > 1 {
		return utils.WriteJson(c, 400, fmt.Sprintf("only accept 1 file ,got '%d'", len(files)), nil)
	}

	if errs := h.XValidator.Validate(&payload); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)
		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf("[%s:%v] need to implement '%s'", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "failed field on :", errMsg)
	}

	file := files[0]

	fileExt := filepath.Ext(file.Filename)

	dir, err := os.Getwd()
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	newUrl := strings.ReplaceAll(payload.Url, " ", "_")

	file.Filename = newUrl + fileExt
	dest := filepath.Join(dir, "img", file.Filename)

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return utils.WriteJson(c, 500, "mkdir error:"+err.Error(), nil)
	}

	if err := c.SaveFile(file, dest); err != nil {
		return utils.WriteJson(c, 500, "savefile err:"+err.Error(), nil)
	}

	payload.Url = file.Filename

	response, err := h.Service.UpdateUserProfile(c.Context(), &payload)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success update profile", response)
}
