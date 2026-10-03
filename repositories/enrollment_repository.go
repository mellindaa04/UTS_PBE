package repositories

import (
	"context"
	"errors"
	"fmt"

	"UTS_PBE/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnrollmentRepository struct {
	DB *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) *EnrollmentRepository {
	return &EnrollmentRepository{
		DB: db,
	}
}

// Mengambil seluruh enrollment milik mahasiswa tertentu.
func (r *EnrollmentRepository) FindByStudentID(
	studentID int,
) ([]models.Enrollment, error) {
	rows, err := r.DB.Query(
		context.Background(),
		`SELECT id, student_id, course_id, tahun_akademik, created_at
		 FROM enrollments
		 WHERE student_id = $1
		 ORDER BY created_at DESC`,
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []models.Enrollment

	for rows.Next() {
		var enrollment models.Enrollment

		err := rows.Scan(
			&enrollment.ID,
			&enrollment.StudentID,
			&enrollment.CourseID,
			&enrollment.TahunAkademik,
			&enrollment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enrollments, nil
}

// Menambahkan enrollment baru.
func (r *EnrollmentRepository) Create(
	studentID int,
	courseID int,
	tahunAkademik string,
) (*models.Enrollment, error) {
	var enrollment models.Enrollment

	err := r.DB.QueryRow(
		context.Background(),
		`INSERT INTO enrollments (student_id, course_id, tahun_akademik)
		 VALUES ($1, $2, $3)
		 RETURNING id, student_id, course_id, tahun_akademik, created_at`,
		studentID,
		courseID,
		tahunAkademik,
	).Scan(
		&enrollment.ID,
		&enrollment.StudentID,
		&enrollment.CourseID,
		&enrollment.TahunAkademik,
		&enrollment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &enrollment, nil
}

// Menghitung total SKS mahasiswa pada tahun akademik tertentu.
func (r *EnrollmentRepository) GetTotalSKS(
	studentID int,
	tahunAkademik string,
) (int, error) {
	var totalSKS int

	err := r.DB.QueryRow(
		context.Background(),
		`SELECT COALESCE(SUM(c.sks), 0)
		 FROM enrollments e
		 JOIN courses c ON e.course_id = c.id
		 WHERE e.student_id = $1
		   AND e.tahun_akademik = $2`,
		studentID,
		tahunAkademik,
	).Scan(&totalSKS)

	if err != nil {
		return 0, err
	}

	return totalSKS, nil
}

// Membuat enrollment menggunakan transaksi database.
func (r *EnrollmentRepository) CreateWithTransaction(
	studentID int,
	courseID int,
	tahunAkademik string,
) (*models.Enrollment, error) {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var enrollment models.Enrollment

	err = tx.QueryRow(
		ctx,
		`INSERT INTO enrollments (student_id, course_id, tahun_akademik)
		 VALUES ($1, $2, $3)
		 RETURNING id, student_id, course_id, tahun_akademik, created_at`,
		studentID,
		courseID,
		tahunAkademik,
	).Scan(
		&enrollment.ID,
		&enrollment.StudentID,
		&enrollment.CourseID,
		&enrollment.TahunAkademik,
		&enrollment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &enrollment, nil
}

// Memeriksa apakah mahasiswa sudah mengambil mata kuliah pada tahun akademik tertentu.
func (r *EnrollmentRepository) Exists(
	studentID int,
	courseID int,
	tahunAkademik string,
) (bool, error) {
	var exists bool

	err := r.DB.QueryRow(
		context.Background(),
		`SELECT EXISTS (
			SELECT 1
			FROM enrollments
			WHERE student_id = $1
			  AND course_id = $2
			  AND tahun_akademik = $3
		)`,
		studentID,
		courseID,
		tahunAkademik,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *EnrollmentRepository) CreateValidated(
	studentID int,
	courseID int,
	tahunAkademik string,
	maxSKS int,
) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Kunci data mahasiswa.
	var lockedStudentID int

	err = tx.QueryRow(ctx,
		`SELECT id
     FROM students
     WHERE id = $1
       AND deleted_at IS NULL
     FOR UPDATE`,
		studentID,
	).Scan(&lockedStudentID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("mahasiswa tidak ditemukan")
		}
		return err
	}

	// Kunci data mata kuliah dan ambil SKS serta kuota.
	var courseSKS, kuota int
	err = tx.QueryRow(ctx,
		`SELECT sks, kuota FROM courses
         WHERE id = $1
         FOR UPDATE`,
		courseID,
	).Scan(&courseSKS, &kuota)
	if err != nil {
		return err
	}

	// Cek duplikasi mata kuliah.
	var duplicate bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS (
            SELECT 1 FROM enrollments
            WHERE student_id = $1
              AND course_id = $2
              AND tahun_akademik = $3
        )`,
		studentID, courseID, tahunAkademik,
	).Scan(&duplicate)
	if err != nil {
		return err
	}
	if duplicate {
		return fmt.Errorf("mata kuliah sudah diambil pada tahun akademik ini")
	}

	// Cek total SKS mahasiswa.
	var totalSKS int
	err = tx.QueryRow(ctx,
		`SELECT COALESCE(SUM(c.sks), 0)
         FROM enrollments e
         JOIN courses c ON c.id = e.course_id
         WHERE e.student_id = $1
           AND e.tahun_akademik = $2`,
		studentID, tahunAkademik,
	).Scan(&totalSKS)
	if err != nil {
		return err
	}

	if totalSKS+courseSKS > maxSKS {
		return fmt.Errorf("batas SKS terlampaui: maksimal %d SKS", maxSKS)
	}

	// Cek kuota kelas.
	var enrolledCount int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM enrollments
         WHERE course_id = $1
           AND tahun_akademik = $2`,
		courseID, tahunAkademik,
	).Scan(&enrolledCount)
	if err != nil {
		return err
	}

	if enrolledCount >= kuota {
		return fmt.Errorf("kuota mata kuliah sudah penuh")
	}

	// Simpan enrollment.
	_, err = tx.Exec(ctx,
		`INSERT INTO enrollments
         (student_id, course_id, tahun_akademik)
         VALUES ($1, $2, $3)`,
		studentID, courseID, tahunAkademik,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Menghapus enrollment milik mahasiswa tertentu.
func (r *EnrollmentRepository) DeleteByID(
	enrollmentID int,
	studentID int,
) error {
	ctx := context.Background()

	commandTag, err := r.DB.Exec(
		ctx,
		`DELETE FROM enrollments
         WHERE id = $1 AND student_id = $2`,
		enrollmentID,
		studentID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("KRS tidak ditemukan atau bukan milik Anda")
	}

	return nil
}
