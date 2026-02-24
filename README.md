# Gowir (E-Commerce Backend)

Ini adalah backend e-commerce yang menggunakan Go, PostgreSQL, `golang-migrate` untuk migrasi database, dan `sqlc` untuk generate query Go yang *type-safe*.

## 🛠️ Stack Teknologi

- **Bahasa**: Go (Golang)
- **Database**: PostgreSQL (via Docker/Podman)
- **Database GUI**: NocoDB
- **Database Driver**: `jackc/pgx/v5`
- **Migration Tool**: `golang-migrate/migrate`
- **Query Builder**: `sqlc`
- **UUID Generator**: `google/uuid` (UUID v7)

## 🚀 Menjalankan Database

Untuk menjalankan PostgreSQL dan NocoDB secara lokal:

```bash
podman-compose up -d
# atau
docker compose up -d
```

- **NocoDB UI:** `http://localhost:${NOCODB_PORT}` *(default: `8080`, lihat `.env`)*
- **Postgres Port:** `5432`

## 🗄️ Database Commands (Taskfile)

Proyek ini menggunakan **Taskfile** (alternatif modern dari Makefile) untuk menyederhanakan perintah-perintah panjang. 

> ⚠️ **Pastikan kamu sudah menginstal Taskfile:** 
> Silakan lihat panduan instalasinya di [https://taskfile.dev/docs/installation](https://taskfile.dev/docs/installation).

Berikut adalah perintah-perintah yang tersedia:

| Perintah | Deskripsi |
| :--- | :--- |
| `task db:up` | Menjalankan PostgreSQL dan NocoDB di background |
| `task db:down` | Mematikan PostgreSQL dan NocoDB |
| `task migrate:create -- nama_tabel` | Membuat file migrasi baru (up & down) |
| `task migrate:up` | Menjalankan migrasi database ke versi terbaru |
| `task migrate:down` | Mengembalikan (rollback) migrasi database 1 langkah |
| `task sqlc:generate` | Meng-generate ulang kode Go dari file SQL menggunakan SQLC |

*(Kamu bisa melihat daftar lengkap ini kapan saja dengan mengetik `task` atau `task -l` di terminal).*

### Mengunduh Dependencies

Jika ada *package* atau module yang error/merah:
```bash
go mod tidy
```

---

## 📁 Arsitektur & Struktur Direktori (Golden Rules)

Proyek ini menggunakan arsitektur **Vertical Slice (Co-location)** yang dikombinasikan dengan pola **1 File 1 Endpoint / Action**. Pendekatan ini dipilih karena paling sesuai dengan filosofi Go: sederhana, mudah dibaca, dan sangat mudah di-maintain.

### 1. Struktur Folder Utama

- **`db/`**: Berisi murni urusan SQL (bukan kode Go).
  - `db/migrations/`: File SQL untuk membuat/merubah tabel (DDL).
  - `db/query/`: File SQL berisi query CRUD yang akan diproses oleh SQLC.
- **`internal/`**: Kode aplikasi utama (private, tidak bisa di-import proyek luar).
  - `internal/db/`: Kode Go hasil *generate* otomatis dari SQLC.
  - `internal/features/`: **Jantung aplikasi.** Semua logika bisnis dibagi berdasarkan fitur (domain), bukan berdasarkan lapisan teknis (seperti controller/service/repo).

### 2. Pola "1 File 1 Endpoint" dalam Fitur

Di dalam setiap folder fitur (misal `internal/features/category/`), kita **TIDAK** memisahkan kode ke dalam `handler.go`, `service.go`, dan `dto.go` yang besar. 

Sebaliknya, kita memecah file berdasarkan **Aksi (Endpoint)**:

```text
internal/features/category/
├── create.go        # Berisi Request/Response Struct, Logika Validasi, Handler, dan Service khusus untuk "Create"
├── get_detail.go    # Logika murni untuk "Get by ID"
├── update.go        # Logika murni untuk "Update"
├── delete.go        # Logika murni untuk "Delete"
└── routes.go        # Hanya berisi pendaftaran rute (misal: router.Post("/", CreateCategory))
```

**Keuntungan:**
- Jika ada *bug* di fitur Update Kategori, kamu hanya perlu membuka `update.go`. Semua konteks (dari request masuk sampai eksekusi database) ada di satu file tersebut dari atas ke bawah.
- Mencegah *file bloating* (file membengkak ribuan baris) dan meminimalisir *merge conflict* di Git.

### 3. Prinsip: Maintainability > DRY (Don't Repeat Yourself)

Sesuai filosofi Go: *"A little copying is better than a little dependency."*

- **Boleh Redundan:** Sangat dianjurkan membuat struct Request/Response (`DTO`) yang terpisah di setiap file aksi (misal struct `CreateCategoryRequest` di `create.go` dan `UpdateCategoryRequest` di `update.go`), meskipun isinya mirip. Ini mencegah fitur `Update` merusak fitur `Create` ketika kebutuhannya mulai berbeda di masa depan.
- **Hindari "God Object" atau Fungsi Helper Terlalu Global:** Jangan memaksakan membuat satu fungsi database raksasa yang dipakai oleh semua *endpoint* jika pada akhirnya fungsi tersebut membutuhkan terlalu banyak `if/else` untuk menyesuaikan kebutuhan masing-masing *endpoint*.
- **Kapan Harus di-Share?** Hanya pisahkan/share logika bisnis inti yang sangat kompleks (misal: algoritma perhitungan pajak) atau kode infrastruktur murni (logger, middleware auth) ke folder seperti `internal/shared/`.

### 4. Menangani Kompleksitas (Orchestrator)

Untuk fitur kompleks yang melibatkan banyak entitas (seperti `checkout`), buatlah folder fitur khusus (misal `internal/features/checkout/`). Fitur ini bertindak sebagai **Orkestrator**.

- Fitur murni (seperti `cart`, `product`) dilarang saling memodifikasi untuk mencegah *circular dependency*.
- Sementara itu, `checkout` bertugas memanggil `cart`, memanggil `product` (cek stok), dan membungkus semuanya dalam satu **Database Transaction (DB TX)**.