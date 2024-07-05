package utils

import (
	"backend/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type MiddlewareStruct struct {
	App fiber.Router
	*config.Config
}

func Middleware(app fiber.Router, c *config.Config) *MiddlewareStruct {
	return &MiddlewareStruct{app, c}
}

func (m *MiddlewareStruct) Middleware(c *fiber.Ctx) error {
	return c.Next()
}

func (m *MiddlewareStruct) MiddlewareWithJWT(c *fiber.Ctx) error {

	token := c.Cookies("jwt", "invalid")
	if token == "invalid" {
		return WriteJson(c, 401, "unautorized", nil)
	}
	jwtToken, err := validateJWT(token, m.Config)
	if err != nil {
		return WriteJson(c, 401, err.Error(), nil)
	}

	if !jwtToken.Valid {
		return WriteJson(c, 401, "unautorized", nil)
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		userID := claims["id"].(string)
		c.Locals("userID", userID)
	}

	return c.Next()

}

func validateJWT(token string, env *config.Config) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(env.JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}
