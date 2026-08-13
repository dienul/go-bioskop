CREATE DATABASE go_bioskop;

-- Setelah database dibuat, hubungkan PostgreSQL ke database go_bioskop.
-- Kemudian jalankan perintah CREATE TABLE berikut.
CREATE TABLE bioskop (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    lokasi VARCHAR(255) NOT NULL,
    rating DOUBLE PRECISION NOT NULL DEFAULT 0
);
