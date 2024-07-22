package api

import (
	"backend/config"
	"backend/server/db"
	"backend/server/internal/cart"
	"backend/server/internal/orders"
	"backend/server/internal/product"
	"backend/server/internal/user"
	"backend/server/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var Env = config.Env.Main().InitConfig()

type Api struct {
	Addr string
}

func NewApi() *Api {
	return &Api{
		Addr: Env.Port,
	}
}

func (a *Api) Run() error {
	fiberApp := fiber.New(fiber.Config{CaseSensitive: true, EnableSplittingOnParsers: true, ErrorHandler: utils.ErrorHandler})

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: " http://localhost:5173",
		// http://localhost:5173
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,OPTIONS",
	}))

	api := fiberApp.Group("/api")

	middleware := utils.Middleware(api, &Env)
	api.Use(middleware.Middleware)

	db, err := db.NewDatabase(&Env)
	if err != nil {
		return err
	}

	validate := &utils.XValidator{
		Validator: validator.New(),
	}

	userRepo := user.NewRepository(db.DB())
	userService := user.NewService(userRepo, &Env)
	userHandler := user.NewHandler(userService, *validate, *middleware, middleware.App)
	userHandler.RegisterRoute()

	productRepo := product.NewRepository(db.DB())
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService, *validate, middleware, middleware.App)
	productHandler.RegisterRoute()

	cartRepo := cart.NewRepository(db.DB())
	cartService := cart.NewService(cartRepo, productRepo)
	cartHandler := cart.NewHandler(cartService, *validate, *middleware, middleware.App)
	cartHandler.RegisterRoute()

	orderRepo := orders.NewRepository(db.DB())
	orderService := orders.NewService(orderRepo, productRepo)
	orderHandler := orders.NewHandler(orderService, *validate, middleware, middleware.App)
	orderHandler.RegisterRoute()

	return fiberApp.Listen(a.Addr)
}
