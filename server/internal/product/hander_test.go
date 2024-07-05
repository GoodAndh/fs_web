package product

import (
	"backend/config"
	"backend/server/db"
	"backend/server/internal/user"
	"backend/server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const baseURL = "http://localhost:3000/api/"

func TestE2E(t *testing.T) {
	fiberApp, err := setUp()
	if err != nil {
		t.Fail()
	}

	t.Run("get all product ", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, baseURL+"product", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}
		byte, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fail()
		}
		var resp utils.GlobalResponseError
		if err := json.Unmarshal(byte, &resp); err != nil {
			t.Fail()
		}
		fmt.Println("response:", resp)
		assert.Equal(t, 200, resp.Status)

	})

	t.Run("get product by id ,fail if given was string /not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"product/id/1", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}
		byte, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fail()
		}
		var resp utils.GlobalResponseError
		if err := json.Unmarshal(byte, &resp); err != nil {
			t.Fail()
		}
		fmt.Println("response:", resp)
		assert.Equal(t, 200, resp.Status)

	})

	t.Run("get product by name ,fail if given was int /not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"product/search/nasi goreng", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}
		byte, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fail()
		}
		var resp utils.GlobalResponseError
		if err := json.Unmarshal(byte, &resp); err != nil {
			t.Fail()
		}
		fmt.Println("response:", resp)
		assert.Equal(t, 200, resp.Status)

	})

	t.Run("create product should fail if unauthorized,invalid payload", func(t *testing.T) {
		payload := &CreateProductRequest{
			Name:        "",
			Description: "",
			Price:       0,
			Stock:       0,
		}
		bite, err := json.Marshal(&payload)
		if err != nil {
			t.Fail()
		}
		req, err := http.NewRequest(http.MethodPost, baseURL+"product/create", bytes.NewBuffer(bite))
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")

		response, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}
		byte, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fail()
		}
		var resp utils.GlobalResponseError
		if err := json.Unmarshal(byte, &resp); err != nil {
			t.Fail()
		}
		fmt.Println("response:", resp)
		assert.Equal(t, 401, resp.Status)

	})

	t.Run("create product success (need to login)", func(t *testing.T) {

		logReq, err := http.NewRequest(http.MethodPost, baseURL+"signin", strings.NewReader(`{"username":"username","password":"password"}`))
		if err != nil {
			t.Fail()
		}
		logReq.Header.Add("Content-Type", "application/json")

		logResponse, err := fiberApp.Test(logReq)
		if err != nil {
			t.Fail()
		}

		logbyte, err := io.ReadAll(logResponse.Body)
		if err != nil {
			t.Fail()
		}
		var logResp utils.GlobalResponseError
		if err := json.Unmarshal(logbyte, &logResp); err != nil {
			t.Fail()
		}

		fmt.Println("login response:", logResp)
		assert.Equal(t, 200, logResp.Status)

		// take cookie
		cookies := logResponse.Cookies()

		payload := &CreateProductRequest{
			Name:        "nasi goreng",
			Description: "nasi goreng spesial",
			Price:       15000,
			Stock:       100,
		}
		bite, err := json.Marshal(&payload)
		if err != nil {
			t.Fail()
		}
		req, err := http.NewRequest(http.MethodPost, baseURL+"product/create", bytes.NewBuffer(bite))
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		response, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}
		byte, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fail()
		}
		var resp utils.GlobalResponseError
		if err := json.Unmarshal(byte, &resp); err != nil {
			t.Fail()
		}
		fmt.Println("response:", resp)

		assert.Equal(t, 200, resp.Status)
	})

}

var Env = config.Env.Test().InitConfig()

func setUp() (*fiber.App, error) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api")

	middleware := utils.Middleware(api, &Env)

	database, err := db.NewDatabase(&Env)
	if err != nil {
		return nil, err
	}

	validate := &utils.XValidator{
		Validator: validator.New(),
	}

	userRepo := user.NewRepository(database.DB())
	userService := user.NewService(userRepo, &Env)
	userHandler := user.NewHandler(userService, *validate)
	userHandler.RegisterRoute(middleware.App)

	repo := NewRepository(database.DB())
	service := NewService(repo)
	handler := NewHandler(service, *validate, middleware)
	handler.RegisterRoute(middleware.App)

	return fiberApp, nil
}
