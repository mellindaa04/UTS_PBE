package services

import (
	"UTS_PBE/repositories"
)

type EnrollmentService struct {
	EnrollmentRepo *repositories.EnrollmentRepository
	StudentRepo    *repositories.StudentRepository
	CourseRepo     *repositories.CourseRepository
}

func NewEnrollmentService(
	enrollmentRepo *repositories.EnrollmentRepository,
	studentRepo *repositories.StudentRepository,
	courseRepo *repositories.CourseRepository,
) *EnrollmentService {
	return &EnrollmentService{
		EnrollmentRepo: enrollmentRepo,
		StudentRepo:    studentRepo,
		CourseRepo:     courseRepo,
	}
}

// Menentukan batas maksimal SKS berdasarkan IPK.
func maxSKSByIPK(ipk float64) int {
	if ipk >= 3.00 {
		return 24
	}

	if ipk >= 2.50 {
		return 21
	}

	return 18
}

func (s *EnrollmentService) Enroll(
	userID int,
	courseID int,
	tahunAkademik string,
) error {
	// Cari data mahasiswa berdasarkan user yang login.
	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	// Tentukan batas SKS berdasarkan IPK.
	maxSKS := maxSKSByIPK(student.IPKTerakhir)

	// Jalankan validasi dan simpan KRS melalui transaksi.
	return s.EnrollmentRepo.CreateValidated(
		student.ID,
		courseID,
		tahunAkademik,
		maxSKS,
	)
}

func (s *EnrollmentService) DeleteEnrollment(
	userID int,
	enrollmentID int,
) error {
	// Cari data mahasiswa berdasarkan akun yang login.
	student, err := s.StudentRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	// Hapus KRS hanya jika dimiliki mahasiswa tersebut.
	return s.EnrollmentRepo.DeleteByID(
		enrollmentID,
		student.ID,
	)
}
