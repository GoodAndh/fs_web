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
	fiber.Router
}

func NewHandler(service Service, validate utils.XValidator, m utils.MiddlewareStruct, r fiber.Router) *Handler {
	return &Handler{service, validate, m, r}
}

func (h *Handler) RegisterRoute() {
	h.Router.Post("/signup", h.createUsers)
	h.Router.Post("/signin", h.loginUsers)
	h.Router.Get("/signout", h.logoutUsers)
	h.Router.Get("/user/profile", h.MiddlewareWithJWT, h.getProfile)
	h.Router.Post("/user/profile/create", h.MiddlewareWithJWT, h.createProfile)
	h.Router.Post("/user/profile/update", h.MiddlewareWithJWT, h.updateProfile)
	h.Router.Get("/user/serveprofile", h.MiddlewareWithJWT, h.serveFile)
	h.Router.Get("/user/validate", h.MiddlewareWithJWT, h.validateJWT)
	h.Router.Get("/user/getuser", h.MiddlewareWithJWT, h.getUsers)
}

func (h *Handler) createUsers(c *fiber.Ctx) error {
	var Req RegisUserRequest
	if err := c.BodyParser(&Req); err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	errs := h.XValidator.Validate(&Req)
	if len(errs) > 0 {
		return utils.WriteJson(c, 400, "failed field on :", errs)
	}

	response, err := h.Service.CreateUsers(c.Context(), &Req)
	if err != nil {
		if err == ErrUsernameAlrInUsed {
			return utils.WriteJson(c, 400, "failed on:", fiber.Map{
				"Username": ErrUsernameAlrInUsed.Error(),
			})
		}
		if err == ErrEmailAlrInUsed {
			return utils.WriteJson(c, 400, "failed on:", fiber.Map{
				"Email": ErrEmailAlrInUsed.Error(),
			})
		}
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
		return utils.WriteJson(c, 400, ErrMissingFile.Error(), fiber.Map{
			"fileMessage": ErrMissingFile.Error(),
		})
	}

	if len(files) > 1 {
		return utils.WriteJson(c, 400, fmt.Sprintf("only accept 1 file ,got '%d'", len(files)), nil)
	}

	errs := h.XValidator.Validate(&payload)
	if len(errs) > 0 {

		return utils.WriteJson(c, 400, "failed field on :", errs)
	}

	file := files[0]

	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return utils.WriteJson(c, 400, "only  image files are allowed", nil)
	}

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
		return utils.WriteJson(c, 400, ErrMissingFile.Error(), fiber.Map{
			"fileMessage": ErrMissingFile.Error(),
		})
	}

	if len(files) > 1 {
		return utils.WriteJson(c, 400, fmt.Sprintf("only accept 1 file ,got '%d'", len(files)), nil)
	}

	errs := h.XValidator.Validate(&payload)
	if len(errs) > 0 {

		return utils.WriteJson(c, 400, "failed field on :", errs)
	}

	file := files[0]
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return utils.WriteJson(c, 400, "only  image files are allowed", nil)
	}

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

func (h *Handler) serveFile(c *fiber.Ctx) error {

	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	profile, err := h.Service.GetUserProfile(c.Context(), userID)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	dir, err := os.Getwd()
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	return c.SendFile(dir + "/img/" + profile.Url)
}

func (h *Handler) validateJWT(c *fiber.Ctx) error {

	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "your token is still available", map[string]int{
		"userID": userID,
	})

}

func (h *Handler) getUsers(c *fiber.Ctx) error {

	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	response, err := h.Service.GetUserByID(c.Context(), userID)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "response success", response)

}
