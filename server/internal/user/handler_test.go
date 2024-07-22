package user

import (
	"backend/config"
	"backend/server/db"
	"backend/server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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

	t.Run("Register User Fail", func(t *testing.T) {
		payload := &RegisUserRequest{
			Username:  "",
			Email:     "",
			Password:  "",
			VPassword: "",
		}

		byte, err := json.Marshal(payload)
		if err != nil {
			t.Fail()
		}

		r, err := http.NewRequest(http.MethodPost, baseURL+"signup", bytes.NewBuffer(byte))
		if err != nil {
			t.Fail()
		}
		r.Header.Add("Content-Type", "application/json")

		response, err := fiberApp.Test(r)
		if err != nil {
			t.Fail()
		}

		bite, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fail()
		}

		var resp utils.GlobalResponseError

		if err := json.Unmarshal(bite, &resp); err != nil {
			t.Fail()
		}

		assert.Equal(t, 400, resp.Status)

	})

	t.Run("Regis User Success", func(t *testing.T) {
		payload := RegisUserRequest{
			Username:  "username",
			Email:     "email@gmail.com",
			Password:  "password",
			VPassword: "password",
		}
		byte, err := json.Marshal(payload)
		if err != nil {
			t.Fail()
		}

		req, err := http.NewRequest(http.MethodPost, baseURL+"signup", bytes.NewBuffer(byte))
		if err != nil {
			t.Fail()
		}

		req.Header.Add("Content-Type", "application/json")
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError
		if err := json.Unmarshal(bytes, &response); err != nil {
			t.Fail()
		}

		if response.Message == "username already in used" || response.Message == "email already in used" {
			assert.Equal(t, 400, response.Status)
			return
		}

		assert.Equal(t, 200, response.Status)

	})

	t.Run("login user failed", func(t *testing.T) {
		payload := &LoginUserRequest{
			Username: "",
			Password: "",
		}

		byte, err := json.Marshal(payload)
		if err != nil {
			t.Fail()
		}

		req, err := http.NewRequest(http.MethodPost, baseURL+"signin", bytes.NewBuffer(byte))
		if err != nil {
			t.Fail()
		}

		req.Header.Add("Content-Type", "application/json")

		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}

		bites, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError

		if err := json.Unmarshal(bites, &response); err != nil {
			t.Fail()
		}

		assert.Equal(t, 400, response.Status)

	})

	t.Run("login user success", func(t *testing.T) {
		payload := &LoginUserRequest{
			Username: "username",
			Password: "password",
		}

		byte, err := json.Marshal(payload)
		if err != nil {
			t.Fail()
		}

		req, err := http.NewRequest(http.MethodPost, baseURL+"signin", bytes.NewBuffer(byte))
		if err != nil {
			t.Fail()
		}

		req.Header.Add("Content-Type", "application/json")

		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}

		bites, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fail()
		}

		var response utils.GlobalResponseError

		if err := json.Unmarshal(bites, &response); err != nil {
			t.Fail()
		}

		assert.Equal(t, 200, response.Status)
	})

	t.Run("create image success", func(t *testing.T) {

		payload := &LoginUserRequest{
			Username: "username",
			Password: "password",
		}

		byte, err := json.Marshal(payload)
		if err != nil {
			t.Fail()
		}

		reqLogin, err := http.NewRequest(http.MethodPost, baseURL+"signin", bytes.NewBuffer(byte))
		if err != nil {
			t.Fail()
		}

		reqLogin.Header.Add("Content-Type", "application/json")

		respLogin, err := fiberApp.Test(reqLogin)
		if err != nil {
			t.Fail()
		}

		cookies := respLogin.Cookies()

		var buf bytes.Buffer

		writer := multipart.NewWriter(&buf)

		part, err := writer.CreateFormFile("file", "file_test.png")
		if err != nil {
			t.Fail()
		}

		file, err := os.Open("file_test.png")
		if err != nil {
			t.Fail()
		}
		defer file.Close()

		_, err = io.Copy(part, file)
		if err != nil {
			t.Fail()
		}
		writer.Close()

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%suser/profile?url=%s&captions=%s", baseURL, "profilebro", "noCAptions"), &buf)
		if err != nil {
			t.Fail()
		}
		req.Header.Add("Content-Type", writer.FormDataContentType())
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fail()
		}
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fail()
		}
		var response utils.GlobalResponseError
		if err := json.Unmarshal(bytes, &response); err != nil {
			t.Fail()
		}
		if response.Message == "you already have profile .do you mean update?" {
			assert.Equal(t, 400, response.Status)
			return
		}
		assert.Equal(t, 200, response.Status)

	})

}

func setUp() (*fiber.App, error) {
	fiberApp := *fiber.New()
	api := fiberApp.Group("/api")

	db, err := db.NewDatabase(&Env)
	if err != nil {
		return nil, err
	}
	middleware := utils.Middleware(api, &Env)

	validate := &utils.XValidator{
		Validator: validator.New(),
	}

	userRepo := NewRepository(db.DB())
	userService := NewService(userRepo, &Env)
	userHandler := NewHandler(userService, *validate, *middleware,middleware.App)
	userHandler.RegisterRoute()

	return &fiberApp, nil
}
