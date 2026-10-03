package handlers

import (
	"UTS_PBE/services"

	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	CourseService *services.CourseService
}

func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		CourseService: courseService,
	}
}

func (h *CourseHandler) GetAll(c *fiber.Ctx) error {
	courses, err := h.CourseService.GetAllCourses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "gagal mengambil data mata kuliah",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": courses,
	})
}
