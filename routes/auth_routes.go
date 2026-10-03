package routes

import (
	"UTS_PBE/handlers"
	"UTS_PBE/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(
	api fiber.Router,
	authHandler *handlers.AuthHandler,
	jwtSecret string,
) {
	auth := api.Group("/auth")

	auth.Post("/login", authHandler.Login)

	auth.Get("/me", middleware.JWTAuth(jwtSecret), authHandler.Me)
}
