package product

import (
	"backend/server/utils"
	"fmt"
	"os"
	"path/filepath"
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
	router.Post("/product/image/:id", h.MiddlewareWithJWT, h.createPIMG)
	router.Get("/product/image/:id", h.MiddlewareWithJWT, h.getImage)
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
	newName := strings.ReplaceAll(name, "%20", " ")
	product, err := h.Service.GetProductByName(c.Context(), newName)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.WriteJson(c, 400, fmt.Sprintf("[name:%v] not found ", newName), nil)
		}
		return utils.WriteJson(c, 400, err.Error(), nil)
	}
	return utils.WriteJson(c, 200, "product by id", product)
}

func (h *Handler) createPIMG(c *fiber.Ctx) error {
	idFromContext, ok := c.Locals("userID").(string)
	if idFromContext == "" || !ok {
		return utils.WriteJson(c, 401, fmt.Sprintf("[id from context :%v ]", idFromContext), nil)
	}

	userID, err := strconv.Atoi(string(idFromContext))
	if err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	prID, err := c.ParamsInt("id", 0)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	url := c.Query("url", "")
	captions := c.Query("captions", "")

	payload := CreateProductImageRequest{
		ProductID: prID,
		Url:       url,
		Captions:  captions,
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

	response, err := h.Service.CreateProductImage(c.Context(), userID, &payload)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success create image", response)
}

func (h *Handler) getImage(c *fiber.Ctx) error {
	prID, err := c.ParamsInt("id", 0)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	response, err := h.Service.GetProductImage(c.Context(), prID)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, fmt.Sprintf("product image of product id '%d'", prID), response)
}
