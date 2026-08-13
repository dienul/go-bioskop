package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Bioskop adalah bentuk data bioskop yang disimpan di database.
type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama" binding:"required"`
	Lokasi string  `json:"lokasi" binding:"required"`
	Rating float64 `json:"rating"`
}

var db *sql.DB

func main() {
	// Overload membaca konfigurasi dari file .env di root project.
	// Nilai dari file .env akan dipakai meskipun terminal memiliki nilai lama.
	if err := godotenv.Overload(); err != nil {
		log.Fatal("Gagal membaca file .env: ", err)
	}

	var err error
	db, err = bukaDatabase()
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}
	defer db.Close()

	router := gin.Default()
	router.POST("/bioskop", tambahBioskop)

	fmt.Println("Server berjalan di http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}

// bukaDatabase membuat koneksi ke PostgreSQL menggunakan environment variable.
func bukaDatabase() (*sql.DB, error) {
	host := ambilEnv("DB_HOST", "localhost")
	port := ambilEnv("DB_PORT", "5432")
	user := ambilEnv("DB_USER", "postgres")
	password := ambilEnv("DB_PASSWORD", "postgres")
	dbName := ambilEnv("DB_NAME", "go_bioskop")
	sslMode := ambilEnv("DB_SSLMODE", "disable")

	koneksi := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode,
	)

	database, err := sql.Open("postgres", koneksi)
	if err != nil {
		return nil, err
	}

	// Ping digunakan untuk memastikan database benar-benar dapat dihubungi.
	if err := database.Ping(); err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}

func ambilEnv(nama, nilaiDefault string) string {
	nilai := os.Getenv(nama)
	if nilai == "" {
		return nilaiDefault
	}
	return nilai
}

// tambahBioskop menerima JSON, memvalidasi data, lalu menyimpannya ke database.
func tambahBioskop(c *gin.Context) {
	var bioskop Bioskop

	if err := c.ShouldBindJSON(&bioskop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Nama dan lokasi tidak boleh kosong",
		})
		return
	}

	query := `
		INSERT INTO bioskop (nama, lokasi, rating)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := db.QueryRow(query, bioskop.Nama, bioskop.Lokasi, bioskop.Rating).Scan(&bioskop.ID)
	if err != nil {
		log.Println("Gagal menjalankan INSERT:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menambahkan bioskop",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bioskop berhasil ditambahkan",
		"data":    bioskop,
	})
}
