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
	fiberApp := fiber.New(fiber.Config{CaseSensitive: true, EnableSplittingOnParsers: true, ErrorHandler: utils.PanicHandler})
	api := fiberApp.Group("/api")
	middleware := utils.Middleware(api, &Env)
	api.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	db, err := db.NewDatabase(&Env)
	if err != nil {
		return err
	}

	validate := &utils.XValidator{
		Validator: validator.New(),
	}

	userRepo := user.NewRepository(db.DB())
	userService := user.NewService(userRepo, &Env)
	userHandler := user.NewHandler(userService, *validate, *middleware)
	userHandler.RegisterRoute(middleware.App)

	productRepo := product.NewRepository(db.DB())
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService, *validate, middleware)
	productHandler.RegisterRoute(middleware.App)

	cartRepo := cart.NewRepository(db.DB())
	cartService := cart.NewService(cartRepo, productRepo)
	cartHandler := cart.NewHandler(cartService, *validate, *middleware)
	cartHandler.RegisterRoute(middleware.App)

	orderRepo := orders.NewRepository(db.DB())
	orderService := orders.NewService(orderRepo, productRepo)
	orderHandler := orders.NewHandler(orderService, *validate, middleware)
	orderHandler.RegisterRoute(middleware.App)

	return fiberApp.Listen(a.Addr)
}
