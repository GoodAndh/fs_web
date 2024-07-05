package cart

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

const baseURL = "http://localhost:3000/api/"

func TestE2E(t *testing.T) {

	fiberApp, err := setUp()
	if err != nil {
		t.Fail()
	}

	t.Run("get cart fail,only fail if token invalid or cart is empy", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"cart", nil)
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

		assert.Equal(t, respTest(resp), true)

	})

	t.Run("get cart success (need to login)", func(t *testing.T) {

		logReq, err := http.NewRequest(http.MethodPost, baseURL+"signin", strings.NewReader(`{"username":"username1","password":"password"}`))
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

		cookies := logResponse.Cookies()

		cartReq, err := http.NewRequest(http.MethodGet, baseURL+"cart", nil)
		if err != nil {
			t.Fail()
		}
		cartReq.Header.Add("Content-Type", "application/json")
		for _, cookie := range cookies {
			cartReq.AddCookie(cookie)
		}

		cartResponse, err := fiberApp.Test(cartReq)
		if err != nil {
			t.Fail()
		}
		byte, err := io.ReadAll(cartResponse.Body)
		if err != nil {
			t.Fail()
		}

		var cartResp utils.GlobalResponseError
		if err := json.Unmarshal(byte, &cartResp); err != nil {
			t.Fail()
		}

		assert.Equal(t, 200, cartResp.Status)

	})

	t.Run("create new cart fail,only fail if unauthorized,invalid payload and alrd have the same cart with the same product", func(t *testing.T) {
		logReq, err := http.NewRequest(http.MethodPost, baseURL+"signin", strings.NewReader(`{"username":"username1","password":"password"}`))
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

		cookies := logResponse.Cookies()

		payload := CreateCartRequest{
			ProductID: 1,
			Status:    "wait",
			Total:     15,
		}

		pByte, err := json.Marshal(&payload)
		if err != nil {
			t.Fail()
		}

		cartReq, err := http.NewRequest(http.MethodPost, baseURL+"cart/create/", bytes.NewBuffer(pByte))
		if err != nil {
			t.Fail()
		}
		cartReq.Header.Add("Content-Type", "application/json")

		for _, cookie := range cookies {
			cartReq.AddCookie(cookie)
		}

		cartRes, err := fiberApp.Test(cartReq)
		if err != nil {
			t.Fail()
		}

		cByte, err := io.ReadAll(cartRes.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError
		if err := json.Unmarshal(cByte, &response); err != nil {
			t.Fail()
		}

		assert.Equal(t, respTest(response), true)

	})

	t.Run("create cart success", func(t *testing.T) {
		logReq, err := http.NewRequest(http.MethodPost, baseURL+"signin", strings.NewReader(`{"username":"username1","password":"password"}`))
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

		cookies := logResponse.Cookies()

		payload := CreateCartRequest{
			ProductID: 2,
			Status:    "wait",
			Total:     15,
		}

		pByte, err := json.Marshal(&payload)
		if err != nil {
			t.Fail()
		}

		cartReq, err := http.NewRequest(http.MethodPost, baseURL+"cart/create/", bytes.NewBuffer(pByte))
		if err != nil {
			t.Fail()
		}
		cartReq.Header.Add("Content-Type", "application/json")

		for _, cookie := range cookies {
			cartReq.AddCookie(cookie)
		}

		cartRes, err := fiberApp.Test(cartReq)
		if err != nil {
			t.Fail()
		}

		cByte, err := io.ReadAll(cartRes.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError
		if err := json.Unmarshal(cByte, &response); err != nil {
			t.Fail()
		}

		assert.Equal(t, 200, response.Status)
	})

}

func respTest(resp utils.GlobalResponseError) bool {
	return resp.Status == 400 || resp.Status == 401
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

	prRepo := product.NewRepository(database.DB())

	cartRepo := NewRepository(database.DB())
	cartService := NewService(cartRepo, prRepo)
	NewHandler(cartService, *validate, *middleware).RegisterRoute(middleware.App)

	return fiberApp, nil
}
