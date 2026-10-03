package handlers

import (
	"UTS_PBE/services"

	"github.com/gofiber/fiber/v2"
)

type StudentHandler struct {
	StudentService *services.StudentService
}

func NewStudentHandler(studentService *services.StudentService) *StudentHandler {
	return &StudentHandler{
		StudentService: studentService,
	}
}

func (h *StudentHandler) GetAll(c *fiber.Ctx) error {
	students, err := h.StudentService.GetAllStudents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data mahasiswa",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data mahasiswa berhasil diambil",
		"data":    students,
	})
}

func (h *StudentHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID mahasiswa tidak valid",
		})
	}

	student, err := h.StudentService.GetStudentByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Mahasiswa tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data mahasiswa berhasil diambil",
		"data":    student,
	})
}

type createStudentRequest struct {
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	NIM         string  `json:"nim"`
	Nama        string  `json:"nama"`
	Prodi       string  `json:"prodi"`
	Angkatan    int     `json:"angkatan"`
	IPKTerakhir float64 `json:"ipk_terakhir"`
}

func (h *StudentHandler) Create(c *fiber.Ctx) error {
	var req createStudentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format request tidak valid",
		})
	}

	student, err := h.StudentService.CreateStudent(
		services.CreateStudentInput{
			Email:       req.Email,
			Password:    req.Password,
			NIM:         req.NIM,
			Nama:        req.Nama,
			Prodi:       req.Prodi,
			Angkatan:    req.Angkatan,
			IPKTerakhir: req.IPKTerakhir,
		},
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Mahasiswa berhasil ditambahkan",
		"data":    student,
	})
}

func (h *StudentHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID mahasiswa tidak valid",
		})
	}

	var req createStudentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format request tidak valid",
		})
	}

	student, err := h.StudentService.UpdateStudent(
		id,
		services.CreateStudentInput{
			NIM:         req.NIM,
			Nama:        req.Nama,
			Prodi:       req.Prodi,
			Angkatan:    req.Angkatan,
			IPKTerakhir: req.IPKTerakhir,
		},
	)
	if err != nil {
		if err.Error() == "mahasiswa tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data mahasiswa berhasil diperbarui",
		"data":    student,
	})
}

func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID mahasiswa tidak valid",
		})
	}

	err = h.StudentService.DeleteStudent(id)
	if err != nil {
		if err.Error() == "mahasiswa tidak ditemukan" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus mahasiswa",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Mahasiswa berhasil dihapus",
	})
}
