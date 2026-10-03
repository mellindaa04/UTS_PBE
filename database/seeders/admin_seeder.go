package seeders

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(pool *pgxpool.Pool) error {
	email := "admin@siakad.com"
	password := "Admin123!"

	// Hash password sebelum disimpan ke database.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	// Masukkan admin jika email belum terdaftar.
	_, err = pool.Exec(
		context.Background(),
		`INSERT INTO users (email, password, role)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (email) DO NOTHING`,
		email,
		string(hashedPassword),
		"admin",
	)
	return err
}
