package repositories

import (
	"context"

	"UTS_PBE/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseRepository struct {
	DB *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{
		DB: db,
	}
}

// Mengambil seluruh data mata kuliah.
func (r *CourseRepository) FindAll() ([]models.Course, error) {
	rows, err := r.DB.Query(
		context.Background(),
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota
		 FROM courses
		 ORDER BY semester, kode_mk`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course

	for rows.Next() {
		var course models.Course

		err := rows.Scan(
			&course.ID,
			&course.KodeMK,
			&course.NamaMK,
			&course.SKS,
			&course.Semester,
			&course.Kuota,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// Mengambil data mata kuliah berdasarkan ID.
func (r *CourseRepository) FindByID(id int) (*models.Course, error) {
	var course models.Course

	err := r.DB.QueryRow(
		context.Background(),
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota
		 FROM courses
		 WHERE id = $1`,
		id,
	).Scan(
		&course.ID,
		&course.KodeMK,
		&course.NamaMK,
		&course.SKS,
		&course.Semester,
		&course.Kuota,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

// Menghitung jumlah mahasiswa yang mengambil mata kuliah pada tahun akademik tertentu.
func (r *CourseRepository) GetEnrolledCount(
	courseID int,
	tahunAkademik string,
) (int, error) {
	var count int

	err := r.DB.QueryRow(
		context.Background(),
		`SELECT COUNT(*)
		 FROM enrollments
		 WHERE course_id = $1
		   AND tahun_akademik = $2`,
		courseID,
		tahunAkademik,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
