package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "app-hexagonaldocs"
	"app-hexagonalinternal/delivery/http"
)

type RouteConfig struct {
	App         *fiber.App
	UserHandler *http.UserHandler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World")
	})

	c.App.Get("/swagger/*", swagger.HandlerDefault) // default swagger
}

func (c *RouteConfig) SetupAuthRoute() {

}
