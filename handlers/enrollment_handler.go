package handlers

import (
	"strconv"
	"strings"

	"UTS_PBE/services"

	"github.com/gofiber/fiber/v2"
)

type EnrollmentHandler struct {
	EnrollmentService *services.EnrollmentService
}

func NewEnrollmentHandler(
	enrollmentService *services.EnrollmentService,
) *EnrollmentHandler {
	return &EnrollmentHandler{
		EnrollmentService: enrollmentService,
	}
}

type CreateEnrollmentRequest struct {
	CourseID      int    `json:"course_id"`
	TahunAkademik string `json:"tahun_akademik"`
}

func (h *EnrollmentHandler) Create(c *fiber.Ctx) error {
	var req CreateEnrollmentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "format request tidak valid",
		})
	}

	req.TahunAkademik = strings.TrimSpace(req.TahunAkademik)

	if req.CourseID <= 0 || req.TahunAkademik == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "course_id dan tahun_akademik wajib diisi",
		})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "identitas pengguna tidak valid",
		})
	}

	err := h.EnrollmentService.Enroll(
		userID,
		req.CourseID,
		req.TahunAkademik,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "KRS berhasil ditambahkan",
	})
}

func (h *EnrollmentHandler) Delete(c *fiber.Ctx) error {
	// Ambil ID enrollment dari URL.
	enrollmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil || enrollmentID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID KRS tidak valid",
		})
	}

	// Ambil ID pengguna dari JWT.
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "identitas pengguna tidak valid",
		})
	}

	// Hapus KRS melalui service.
	err = h.EnrollmentService.DeleteEnrollment(userID, enrollmentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "KRS berhasil dihapus",
	})
}
