package main

import (
	"log"

	"github.com/joho/godotenv"

	"UTS_PBE/config"
	"UTS_PBE/database/seeders"
)

func main() {
	// Memuat konfigurasi dari file .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env:", err)
	}

	// Koneksi ke database
	pool, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}
	defer pool.Close()

	// Seeder admin
	if err := seeders.SeedAdmin(pool); err != nil {
		log.Fatal("Gagal melakukan seeder admin:", err)
	}
	log.Println("Seeder admin berhasil")

	// Seeder mahasiswa
	if err := seeders.SeedStudents(pool); err != nil {
		log.Fatal("Gagal melakukan seeder mahasiswa:", err)
	}
	log.Println("Seeder mahasiswa berhasil")

	// Seeder mata kuliah
	if err := seeders.SeedCourses(pool); err != nil {
		log.Fatal("Gagal melakukan seeder mata kuliah:", err)
	}
	log.Println("Seeder mata kuliah berhasil")

	log.Println("Semua seeder berhasil dijalankan")
}
