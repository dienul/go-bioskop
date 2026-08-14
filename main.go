package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	router.GET("/bioskop", ambilSemuaBioskop)
	router.GET("/bioskop/:id", ambilBioskop)
	router.PUT("/bioskop/:id", ubahBioskop)
	router.DELETE("/bioskop/:id", hapusBioskop)

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

	if err := c.ShouldBindJSON(&bioskop); err != nil || !inputBioskopValid(bioskop) {
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

// ambilSemuaBioskop mengambil seluruh data bioskop dari database.
func ambilSemuaBioskop(c *gin.Context) {
	rows, err := db.Query("SELECT id, nama, lokasi, rating FROM bioskop ORDER BY id")
	if err != nil {
		log.Println("Gagal menjalankan SELECT:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data bioskop"})
		return
	}
	defer rows.Close()

	daftarBioskop := []Bioskop{}
	for rows.Next() {
		var bioskop Bioskop
		if err := rows.Scan(&bioskop.ID, &bioskop.Nama, &bioskop.Lokasi, &bioskop.Rating); err != nil {
			log.Println("Gagal membaca data bioskop:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membaca data bioskop"})
			return
		}
		daftarBioskop = append(daftarBioskop, bioskop)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error saat membaca daftar bioskop:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membaca data bioskop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": daftarBioskop})
}

// ambilBioskop mengambil satu bioskop berdasarkan ID.
func ambilBioskop(c *gin.Context) {
	id, valid := ambilID(c)
	if !valid {
		return
	}

	var bioskop Bioskop
	err := db.QueryRow(
		"SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1",
		id,
	).Scan(&bioskop.ID, &bioskop.Nama, &bioskop.Lokasi, &bioskop.Rating)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bioskop tidak ditemukan"})
		return
	}
	if err != nil {
		log.Println("Gagal mengambil detail bioskop:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil detail bioskop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bioskop})
}

// ubahBioskop memperbarui bioskop berdasarkan ID.
func ubahBioskop(c *gin.Context) {
	id, valid := ambilID(c)
	if !valid {
		return
	}

	var bioskop Bioskop
	if err := c.ShouldBindJSON(&bioskop); err != nil || !inputBioskopValid(bioskop) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama dan lokasi tidak boleh kosong"})
		return
	}

	err := db.QueryRow(`
		UPDATE bioskop
		SET nama = $1, lokasi = $2, rating = $3
		WHERE id = $4
		RETURNING id
	`, bioskop.Nama, bioskop.Lokasi, bioskop.Rating, id).Scan(&bioskop.ID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bioskop tidak ditemukan"})
		return
	}
	if err != nil {
		log.Println("Gagal menjalankan UPDATE:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui bioskop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bioskop berhasil diperbarui",
		"data":    bioskop,
	})
}

// hapusBioskop menghapus bioskop berdasarkan ID.
func hapusBioskop(c *gin.Context) {
	id, valid := ambilID(c)
	if !valid {
		return
	}

	hasil, err := db.Exec("DELETE FROM bioskop WHERE id = $1", id)
	if err != nil {
		log.Println("Gagal menjalankan DELETE:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus bioskop"})
		return
	}

	jumlah, err := hasil.RowsAffected()
	if err != nil {
		log.Println("Gagal membaca jumlah data yang dihapus:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus bioskop"})
		return
	}
	if jumlah == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bioskop tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil dihapus"})
}

// ambilID mengubah parameter ID dari string menjadi integer positif.
func ambilID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID bioskop tidak valid"})
		return 0, false
	}
	return id, true
}

func inputBioskopValid(bioskop Bioskop) bool {
	return strings.TrimSpace(bioskop.Nama) != "" && strings.TrimSpace(bioskop.Lokasi) != ""
}
