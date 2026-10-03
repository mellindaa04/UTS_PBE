package routes

import (
	"UTS_PBE/handlers"
	"UTS_PBE/middleware"

	"github.com/gofiber/fiber/v2"
)

func EnrollmentRoutes(api fiber.Router, enrollmentHandler *handlers.EnrollmentHandler, jwtSecret string) {
	enrollments := api.Group("/enrollments", middleware.JWTAuth(jwtSecret))

	enrollments.Post("/", middleware.RequireRole("student"), enrollmentHandler.Create)
	enrollments.Delete("/:id", middleware.RequireRole("student"), enrollmentHandler.Delete)
}
