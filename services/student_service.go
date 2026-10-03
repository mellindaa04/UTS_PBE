package services

import (
	"errors"
	"strings"

	"UTS_PBE/models"
	"UTS_PBE/repositories"

	"golang.org/x/crypto/bcrypt"
)

type StudentService struct {
	StudentRepo *repositories.StudentRepository
}

func NewStudentService(
	studentRepo *repositories.StudentRepository,
) *StudentService {
	return &StudentService{
		StudentRepo: studentRepo,
	}
}

func (s *StudentService) GetAllStudents() ([]models.Student, error) {
	return s.StudentRepo.FindAll()
}
func (s *StudentService) GetStudentByID(id int) (*models.Student, error) {
	return s.StudentRepo.FindByID(id)
}

func (s *StudentService) UpdateStudent(
	id int,
	input CreateStudentInput,
) (*models.Student, error) {
	input.NIM = strings.TrimSpace(input.NIM)
	input.Nama = strings.TrimSpace(input.Nama)
	input.Prodi = strings.TrimSpace(input.Prodi)

	if input.NIM == "" || input.Nama == "" || input.Prodi == "" {
		return nil, errors.New("semua field wajib diisi")
	}

	if input.Angkatan <= 0 {
		return nil, errors.New("angkatan tidak valid")
	}

	if input.IPKTerakhir < 0 || input.IPKTerakhir > 4 {
		return nil, errors.New("IPK harus antara 0 sampai 4")
	}

	student := models.Student{
		ID:          id,
		NIM:         input.NIM,
		Nama:        input.Nama,
		Prodi:       input.Prodi,
		Angkatan:    input.Angkatan,
		IPKTerakhir: input.IPKTerakhir,
	}

	err := s.StudentRepo.Update(student)
	if err != nil {
		return nil, err
	}

	return s.StudentRepo.FindByID(id)
}

func (s *StudentService) DeleteStudent(id int) error {
	return s.StudentRepo.SoftDelete(id)
}

type CreateStudentInput struct {
	Email       string
	Password    string
	NIM         string
	Nama        string
	Prodi       string
	Angkatan    int
	IPKTerakhir float64
}

func (s *StudentService) CreateStudent(
	input CreateStudentInput,
) (*models.Student, error) {
	input.Email = strings.TrimSpace(input.Email)
	input.NIM = strings.TrimSpace(input.NIM)
	input.Nama = strings.TrimSpace(input.Nama)
	input.Prodi = strings.TrimSpace(input.Prodi)

	if input.Email == "" ||
		input.Password == "" ||
		input.NIM == "" ||
		input.Nama == "" ||
		input.Prodi == "" {
		return nil, errors.New("semua field wajib diisi")
	}

	if len(input.Password) < 8 {
		return nil, errors.New("password minimal 8 karakter")
	}

	if input.Angkatan <= 0 {
		return nil, errors.New("angkatan tidak valid")
	}

	if input.IPKTerakhir < 0 || input.IPKTerakhir > 4 {
		return nil, errors.New("IPK harus antara 0 sampai 4")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	student := models.Student{
		NIM:         input.NIM,
		Nama:        input.Nama,
		Prodi:       input.Prodi,
		Angkatan:    input.Angkatan,
		IPKTerakhir: input.IPKTerakhir,
	}

	return s.StudentRepo.CreateWithUser(
		input.Email,
		string(passwordHash),
		student,
	)
}
