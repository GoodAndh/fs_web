package orders

import (
	"backend/config"
	"backend/server/db"
	"backend/server/internal/product"
	"backend/server/internal/user"
	"backend/server/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	Env = config.Env.Test().InitConfig()
)

const baseURL = "http://localhost:3000/api/"

func TestE2E(t *testing.T) {
	fiberApp, err := setUp()
	if err != nil {
		t.Fail()
	}

	logReq, err := http.NewRequest(http.MethodPost, baseURL+"signin", strings.NewReader(`{"username":"username1","password":"password"}`))
	if err != nil {
		t.Fail()
	}
	logReq.Header.Add("Content-Type", "application/json")

	logResponse, err := fiberApp.Test(logReq)
	if err != nil {
		t.Fail()
	}

	cookies := logResponse.Cookies()
	t.Run("get orders success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"order", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		res, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}

		byte, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fail()
		}
		var response utils.GlobalResponseError
		if err := json.Unmarshal(byte, &response); err != nil {
			t.Fail()
		}
		assert.Equal(t, 200, response.Status)
	})

	t.Run("get order fail unathorized", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"order", nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}

		byte, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError
		if err := json.Unmarshal(byte, &response); err != nil {
			t.Fail()
		}

		assert.Equal(t, 401, response.Status)
	})

	t.Run("create order success", func(t *testing.T) {

		payload := CreateOrderRequest{
			ProductID: 1,
			Total:     100,
		}

		bodyres, err := json.Marshal(&payload)
		if err != nil {
			t.Fail()
		}
		req, err := http.NewRequest(http.MethodPost, baseURL+"order/create", bytes.NewBuffer(bodyres))
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		res, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fail()
		}

		if response.Message == "you alr have the same order ,do you mean change?" {
			assert.Equal(t, 400, response.Status)
			return
		}

		assert.Equal(t, 200, response.Status)

	})

	t.Run("create order fail", func(t *testing.T) {
		payload := CreateOrderRequest{
			ProductID: 1,
			Total:     10000,
		}

		bodyres, err := json.Marshal(&payload)
		if err != nil {
			t.Fail()
		}
		req, err := http.NewRequest(http.MethodPost, baseURL+"order/create", bytes.NewBuffer(bodyres))
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}

		res, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fail()
		}

		

		assert.Equal(t, 400, response.Status)

	})

}

func setUp() (*fiber.App, error) {
	fiberApp := *fiber.New()

	api := fiberApp.Group("/api")
	db, err := db.NewDatabase(&Env)
	if err != nil {
		return nil, err
	}
	middleware := utils.MiddlewareStruct{
		App:    api,
		Config: &Env,
	}

	validate := &utils.XValidator{
		Validator: validator.New(),
	}

	userRepo := user.NewRepository(db.DB())
	userService := user.NewService(userRepo, &Env)
	user.NewHandler(userService, *validate, middleware).RegisterRoute(middleware.App)

	prRepo := product.NewRepository(db.DB())
	orderRepo := NewRepository(db.DB())
	orderService := NewService(orderRepo, prRepo)
	orderHandler := NewHandler(orderService, *validate, &middleware)
	orderHandler.RegisterRoute(middleware.App)

	return &fiberApp, nil

}
