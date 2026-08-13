# Go Bioskop

Aplikasi web sederhana untuk menambahkan data bioskop menggunakan Go, Gin, dan PostgreSQL.

## Struktur data

Data bioskop memiliki atribut berikut:

- `id`: integer dan primary key
- `nama`: string, wajib diisi
- `lokasi`: string, wajib diisi
- `rating`: float

## Cara menjalankan

1. Pastikan Go dan PostgreSQL sudah terpasang.
2. Jalankan isi file `schema.sql` di PostgreSQL.
3. Salin `.env.example` menjadi `.env`, lalu isi `DB_PASSWORD` dengan password PostgreSQL Anda. Aplikasi akan membaca file `.env` secara otomatis.
4. Unduh dependency:

   ```bash
   go mod tidy
   ```

5. Jalankan aplikasi:

   ```bash
   go run .
   ```

Server berjalan di `http://localhost:8080`.

## Mencoba endpoint

Kirim request berikut menggunakan Postman atau aplikasi sejenis:

- Method: `POST`
- URL: `http://localhost:8080/bioskop`
- Body → raw → JSON:

  ```json
  {
    "nama": "Cinema XXI",
    "lokasi": "Jakarta",
    "rating": 4.5
  }
  ```

Jika berhasil, server mengembalikan status `201 Created`:

```json
{
  "message": "Bioskop berhasil ditambahkan",
  "data": {
    "id": 1,
    "nama": "Cinema XXI",
    "lokasi": "Jakarta",
    "rating": 4.5
  }
}
```

Jika `nama` atau `lokasi` tidak diisi, server mengembalikan status `400 Bad Request`.
