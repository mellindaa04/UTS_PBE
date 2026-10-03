package seeders

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseSeed struct {
	KodeMK   string
	NamaMK   string
	SKS      int
	Semester int
	Kuota    int
}

func SeedCourses(pool *pgxpool.Pool) error {
	courses := []CourseSeed{
		{"IF101", "Algoritma dan Pemrograman", 3, 1, 40},
		{"IF102", "Matematika Diskrit", 3, 1, 40},
		{"IF201", "Struktur Data", 3, 2, 35},
		{"IF202", "Basis Data", 3, 2, 35},
		{"IF301", "Pemrograman Web", 3, 3, 30},
		{"IF302", "Rekayasa Perangkat Lunak", 3, 3, 30},
		{"IF401", "Jaringan Komputer", 3, 4, 30},
		{"IF402", "Sistem Operasi", 3, 4, 30},
		{"IF501", "Keamanan Siber", 3, 5, 25},
		{"IF502", "Kecerdasan Buatan", 3, 5, 25},
	}

	ctx := context.Background()

	for _, course := range courses {
		_, err := pool.Exec(ctx,
			`INSERT INTO courses
			 (kode_mk, nama_mk, sks, semester, kuota)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (kode_mk) DO NOTHING`,
			course.KodeMK,
			course.NamaMK,
			course.SKS,
			course.Semester,
			course.Kuota,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
