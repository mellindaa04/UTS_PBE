package repositories

import (
	"context"
	"errors"

	"UTS_PBE/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	DB *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{
		DB: db,
	}
}

// Mengambil seluruh mahasiswa yang belum dihapus.
func (r *StudentRepository) FindAll() ([]models.Student, error) {
	rows, err := r.DB.Query(
		context.Background(),
		`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir
		 FROM students
		 WHERE deleted_at IS NULL
		 ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student

		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.NIM,
			&student.Nama,
			&student.Prodi,
			&student.Angkatan,
			&student.IPKTerakhir,
		)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// Membuat akun pengguna sekaligus profil mahasiswa dalam satu transaksi.
func (r *StudentRepository) CreateWithUser(
	email string,
	passwordHash string,
	student models.Student,
) (*models.Student, error) {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var userID int

	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (email, password, role)
		 VALUES ($1, $2, 'student')
		 RETURNING id`,
		email,
		passwordHash,
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		ctx,
		`INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		userID,
		student.NIM,
		student.Nama,
		student.Prodi,
		student.Angkatan,
		student.IPKTerakhir,
	).Scan(&student.ID)
	if err != nil {
		return nil, err
	}

	student.UserID = userID

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &student, nil
}

// Mengambil mahasiswa berdasarkan ID yang belum dihapus.
func (r *StudentRepository) FindByID(id int) (*models.Student, error) {
	var student models.Student

	err := r.DB.QueryRow(
		context.Background(),
		`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir
		 FROM students
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&student.ID,
		&student.UserID,
		&student.NIM,
		&student.Nama,
		&student.Prodi,
		&student.Angkatan,
		&student.IPKTerakhir,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

// Memperbarui data mahasiswa yang belum dihapus.
func (r *StudentRepository) Update(student models.Student) error {
	commandTag, err := r.DB.Exec(
		context.Background(),
		`UPDATE students
		 SET nim = $1,
		     nama = $2,
		     prodi = $3,
		     angkatan = $4,
		     ipk_terakhir = $5
		 WHERE id = $6 AND deleted_at IS NULL`,
		student.NIM,
		student.Nama,
		student.Prodi,
		student.Angkatan,
		student.IPKTerakhir,
		student.ID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("mahasiswa tidak ditemukan")
	}

	return nil
}

// Menghapus mahasiswa secara soft delete.
func (r *StudentRepository) SoftDelete(id int) error {
	commandTag, err := r.DB.Exec(
		context.Background(),
		`UPDATE students
		 SET deleted_at = NOW()
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("mahasiswa tidak ditemukan")
	}

	return nil
}

// Mengambil mahasiswa berdasarkan user ID yang belum dihapus.
func (r *StudentRepository) FindByUserID(userID int) (*models.Student, error) {
	var student models.Student

	err := r.DB.QueryRow(
		context.Background(),
		`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir
		 FROM students
		 WHERE user_id = $1 AND deleted_at IS NULL`,
		userID,
	).Scan(
		&student.ID,
		&student.UserID,
		&student.NIM,
		&student.Nama,
		&student.Prodi,
		&student.Angkatan,
		&student.IPKTerakhir,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}
