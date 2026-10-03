package seeders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type StudentSeed struct {
	NIM      string
	Nama     string
	Prodi    string
	Angkatan int
	IPK      float64
}

func SeedStudents(pool *pgxpool.Pool) error {
	students := []StudentSeed{
		{"434231001", "Andi Pratama", "Teknik Informatika", 2024, 3.50},
		{"434231002", "Budi Santoso", "Teknik Informatika", 2024, 3.20},
		{"434231003", "Citra Lestari", "Sistem Informasi", 2024, 3.75},
		{"434231004", "Dewi Anggraini", "Teknik Informatika", 2024, 3.10},
		{"434231005", "Eko Saputra", "Sistem Informasi", 2024, 2.90},
		{"434231006", "Fitri Handayani", "Teknik Informatika", 2024, 3.80},
		{"434231007", "Gilang Ramadhan", "Sistem Informasi", 2024, 3.00},
		{"434231008", "Hana Putri", "Teknik Informatika", 2024, 3.40},
		{"434231009", "Indra Wijaya", "Sistem Informasi", 2024, 2.75},
		{"434231010", "Joko Susilo", "Teknik Informatika", 2024, 3.60},
		{"434231011", "Kartika Sari", "Sistem Informasi", 2024, 3.25},
		{"434231012", "Lukman Hakim", "Teknik Informatika", 2024, 2.40},
		{"434231013", "Maya Salsabila", "Sistem Informasi", 2024, 3.90},
		{"434231014", "Nanda Prakoso", "Teknik Informatika", 2024, 3.15},
		{"434231015", "Olivia Maharani", "Sistem Informasi", 2024, 3.55},
		{"434231016", "Putra Firmansyah", "Teknik Informatika", 2024, 2.85},
		{"434231017", "Qori Aulia", "Sistem Informasi", 2024, 3.70},
		{"434231018", "Rizky Kurniawan", "Teknik Informatika", 2024, 3.05},
		{"434231019", "Sinta Amelia", "Sistem Informasi", 2024, 3.35},
		{"434231020", "Taufik Hidayat", "Teknik Informatika", 2024, 2.60},
	}

	ctx := context.Background()

	for _, student := range students {
		email := fmt.Sprintf("%s@siakad.com", student.NIM)

		// Buat akun mahasiswa jika belum ada.
		password, err := bcrypt.GenerateFromPassword(
			[]byte("Mahasiswa123!"),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}

		var userID int
		err = pool.QueryRow(ctx,
			`INSERT INTO users (email, password, role)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (email)
			 DO UPDATE SET email = EXCLUDED.email
			 RETURNING id`,
			email, string(password), "student",
		).Scan(&userID)
		if err != nil {
			return err
		}

		// Buat profil mahasiswa jika belum ada.
		_, err = pool.Exec(ctx,
			`INSERT INTO students
			 (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (nim) DO NOTHING`,
			userID,
			student.NIM,
			student.Nama,
			student.Prodi,
			student.Angkatan,
			student.IPK,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
