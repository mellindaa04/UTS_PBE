package services

import (
	"UTS_PBE/repositories"

	"UTS_PBE/models"
)

type CourseService struct {
	CourseRepo *repositories.CourseRepository
}

func NewCourseService(courseRepo *repositories.CourseRepository) *CourseService {
	return &CourseService{
		CourseRepo: courseRepo,
	}
}

func (s *CourseService) GetAllCourses() ([]models.Course, error) {
	return s.CourseRepo.FindAll()
}
