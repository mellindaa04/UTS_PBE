package routes

import (
	"UTS_PBE/handlers"
	"UTS_PBE/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(
	api fiber.Router,
	studentHandler *handlers.StudentHandler,
	jwtSecret string,
) {
	students := api.Group(
		"/students",
		middleware.JWTAuth(jwtSecret),
	)

	students.Get("/", studentHandler.GetAll)
	students.Get("/:id", studentHandler.GetByID)
	students.Post("/", middleware.RequireRole("admin"), studentHandler.Create)
	students.Put("/:id", middleware.RequireRole("admin"), studentHandler.Update)
	students.Delete("/:id", middleware.RequireRole("admin"), studentHandler.Delete)
}
