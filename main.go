package main

import (
	"log"
	"os"

	"UTS_PBE/config"
	"UTS_PBE/handlers"
	"UTS_PBE/repositories"
	"UTS_PBE/routes"
	"UTS_PBE/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Gagal menghubungkan database:", err)
	}
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET belum dikonfigurasi")
	}

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	studentRepo := repositories.NewStudentRepository(db)
	studentService := services.NewStudentService(studentRepo)
	studentHandler := handlers.NewStudentHandler(studentService)

	courseRepo := repositories.NewCourseRepository(db)
	courseService := services.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseService)

	enrollmentRepo := repositories.NewEnrollmentRepository(db)

	enrollmentService := services.NewEnrollmentService(
		enrollmentRepo,
		studentRepo,
		courseRepo,
	)

	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)

	app := fiber.New()

	api := app.Group("/api/v1")
	routes.RegisterAuthRoutes(api, authHandler, jwtSecret)
	routes.RegisterStudentRoutes(api, studentHandler, jwtSecret)
	routes.CourseRoutes(api, courseHandler, jwtSecret)
	routes.EnrollmentRoutes(api, enrollmentHandler, jwtSecret)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "SIAKAD Mini API berjalan",
		})
	})

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
