package routes

import (
	"UTS_PBE/handlers"
	"UTS_PBE/middleware"

	"github.com/gofiber/fiber/v2"
)

func CourseRoutes(api fiber.Router, courseHandler *handlers.CourseHandler, jwtSecret string) {
	courses := api.Group("/courses", middleware.JWTAuth(jwtSecret))

	courses.Get("/", courseHandler.GetAll)
}
