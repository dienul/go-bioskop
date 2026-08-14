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

## Pengujian dengan curl di PowerShell

Jalankan server dengan `go run .`, lalu buka PowerShell baru untuk menjalankan perintah berikut.

### Tambah bioskop

```powershell
curl.exe -X POST "http://localhost:8080/bioskop" -H "Content-Type: application/json" -d '{\"nama\":\"Cinema XXI\",\"lokasi\":\"Jakarta\",\"rating\":4.5}'
```

### Ambil semua bioskop

```powershell
curl.exe "http://localhost:8080/bioskop"
```

### Ambil detail bioskop

```powershell
curl.exe "http://localhost:8080/bioskop/1"
```

### Perbarui bioskop

```powershell
curl.exe -X PUT "http://localhost:8080/bioskop/1" -H "Content-Type: application/json" -d '{\"nama\":\"Cinema XXI Update\",\"lokasi\":\"Bandung\",\"rating\":4.8}'
```

### Hapus bioskop

```powershell
curl.exe -X DELETE "http://localhost:8080/bioskop/1"
```

Ganti angka `1` dengan ID bioskop yang ingin dilihat, diperbarui, atau dihapus.

## Menggunakan PostgreSQL Railway

Aplikasi akan memakai `DATABASE_URL` jika variabel tersebut tersedia. Konfigurasi `DB_HOST`, `DB_PORT`, dan variabel `DB_*` lainnya hanya menjadi fallback untuk PostgreSQL lokal.

### Menjalankan aplikasi dari komputer lokal

Gunakan URL publik dari TCP Proxy Railway, bukan `RAILWAY_PRIVATE_DOMAIN`. Isi `.env` seperti berikut:

```env
DATABASE_URL=postgresql://postgres:PASSWORD@HOST_PUBLIC:PORT_PUBLIC/railway
PORT=8080
```

Nilai host dan port publik dapat disalin dari halaman Connect PostgreSQL di Railway. Jangan commit file `.env`.

### Menjalankan aplikasi di Railway

Pada service aplikasi Go, tambahkan reference variable yang mengarah ke service PostgreSQL:

```env
DATABASE_URL=${{Postgres.DATABASE_URL}}
```

Ganti `Postgres` jika nama service database Anda berbeda. Railway akan memberikan variabel `PORT` secara otomatis.

Setelah terhubung, buat tabel pada database Railway:

```sql
CREATE TABLE IF NOT EXISTS bioskop (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    lokasi VARCHAR(255) NOT NULL,
    rating DOUBLE PRECISION NOT NULL DEFAULT 0
);
```
