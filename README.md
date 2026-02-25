# POS Warkop

POS kasir sederhana untuk warkop/cafe: POS + cetak struk & kitchen ticket + dashboard finance + setting toko + manajemen user (role).

## Fitur

- POS (Kasir)
  - Cari menu, tambah ke keranjang, atur qty
  - Pembayaran Cash / QRIS
  - Dine In / Take Away
  - Input meja, jumlah tamu, nama pemesan
  - Cetak:
    - Struk customer
    - Orderan kitchen (ticket)
- Setting (Admin)
  - Profil toko: nama, tagline, alamat, telepon
  - Upload logo (dipakai di header + struk)
  - Manajemen user: cashier, finance, admin (superadmin)
- Finance (Admin/Finance)
  - Trace transaksi harian (list order per tanggal + detail)
  - Pembukuan manual (income/expense) dengan metode (cash/bank) dan kategori
  - Neraca harian sederhana (saldo awal kas/bank per tanggal)

## Roles (Akses)

- `cashier`: hanya POS, bisa buat order
- `finance`: hanya Finance, tidak bisa akses POS atau Setting
- `admin`: akses semua (POS + Finance + Setting + Users)

## Subscription & Plan (Standard/Premium)

Project ini mendukung mode subscription sederhana berbasis DB (single-tenant).

- `plan`: `standard` = POS saja (Finance dikunci), `premium` = POS + Finance + Ledger + role finance
- `paid_until`: kalau sudah lewat, sistem masuk mode terbatas (blok `POST /api/orders` dan blok ledger `POST/DELETE /api/ledger`)

Catatan:

- Untuk deployment on-prem di VPS kamu, kamu bisa atur `plan` dan `paid_until` manual (tanpa payment gateway).
- Ini bukan lisensi anti-tamper; user yang punya akses DB tetap bisa mengubah field ini. Untuk model SaaS/multi-tenant perlu mekanisme license server terpisah.

## Tech Stack

- Backend: Go 1.22, Gorilla Mux, Postgres, JWT
- Frontend: Vue 3, Vite, TailwindCSS, Pinia
- Infra: Docker Compose, Postgres init SQL migrations

## Environment Variables

File `.env` dipakai oleh `docker compose` untuk service `db` dan `backend`.

Variabel utama:

- `POSTGRES_USER`: user DB
- `POSTGRES_PASSWORD`: password DB
- `POSTGRES_DB`: nama database
- `DB_HOST`: host DB untuk backend (di docker: `db`)
- `DB_PORT`: port DB (default `5432`)
- `JWT_SECRET`: secret untuk sign JWT (wajib diganti kalau deployment)

## Quick Start (Docker)

### Prasyarat

- Docker Desktop (atau Docker Engine + Compose v2)

### Jalankan

```bash
docker compose up -d --build
```

Buka:

- Frontend UI: `http://localhost:3000`
- Backend API: `http://localhost:8080`

### Akun Seed (Default)

- Email: `admin@pos.com`
- Password: `123456`

### Buat User Kasir / Finance

Cara paling gampang:

1. Login sebagai admin
2. Buka menu `Setting` -> section `User`
3. Klik `Tambah user`, pilih role `Kasir` atau `Finance`
4. Logout, lalu login pakai user baru

### Stop / Reset

Stop containers:

```bash
docker compose down
```

Reset database (hapus semua data):

```bash
docker compose down -v
```

## Database & Migrations

Folder `migrations/` otomatis dijalankan Postgres saat DB pertama kali dibuat (saat volume kosong).

Kalau kamu sudah punya DB lama (volume sudah ada), file migration baru tidak otomatis ter-apply. Gunakan script:

```bash
./scripts/db_apply_migrations.sh
```

Daftar migration:

- `001_init.sql`: tabel utama + seed admin + seed produk dummy
- `002_settings_users_finance.sql`: store settings + role check + order metadata + cash ledger
- `003_store_logo.sql`: kolom logo untuk store settings
- `004_ledger_details.sql`: ledger detail (payment_method, category)
- `005_subscription.sql`: subscription fields (`plan`, `paid_until`)

Set subscription (plan & masa aktif):

```bash
./scripts/subscription_set.sh premium 2026-03-25
./scripts/subscription_set.sh standard 2026-03-25
```

## Setting Toko (Nama/Alamat/Telp/Logo)

Halaman: `/settings` (role admin).

Yang disimpan:

- Nama toko, tagline
- Alamat (2 baris)
- Telepon
- Logo (data URL, max kira-kira 450KB di backend, UI upload dibatasi 250KB)

Dipakai di:

- Header aplikasi
- Struk (opsional, bisa dimatikan dari config receipt)

## Project Structure

- `backend/`: REST API (Go)
- `frontend/`: Web UI (Vue)
- `migrations/`: SQL schema/seed (Postgres init)
- `scripts/`: helper scripts (contoh apply migrations ke DB yang sudah ada)

## Receipt & Kitchen Printing (Thermal)

Config utama:

- `frontend/src/config/receipt.js`

Yang bisa diubah:

- Ukuran kertas printer thermal:
  - 58mm: default `paperWidthMm: 58`, `contentWidthMm: 48`
  - 80mm: umum dipakai `paperWidthMm: 80`, `contentWidthMm: 72`
- Footer struk, pajak, tampilkan logo, tampilkan alamat/telp
- Nama station kitchen: `receiptConfig.kitchen.stationName`

Catatan:

- Print memakai iframe hidden (lebih stabil dibanding popup).
- Page height di-auto-fit biar tidak banyak kertas kosong.
- Saat ada logo di struk, sistem menunggu image load sebentar sebelum memanggil `window.print()`.

## Finance (Pembukuan & Neraca Sederhana)

Halaman: `/finance` (role admin/finance).

Fitur:

- Pilih tanggal, lihat list order pada tanggal itu
- Input pembukuan:
  - `income` / `expense`
  - metode: `cash` / `bank` (anggap bank = transfer/QRIS/dll)
  - kategori: bebas (contoh: `bahan`, `operasional`)
  - keterangan + nominal
- Neraca harian sederhana:
  - Isi saldo awal kas & saldo awal bank/QRIS
  - Sistem hitung saldo akhir kas/bank dan total aktiva vs modal

Catatan:

- Ini "neraca sederhana" untuk kebutuhan operasional warkop.
- Kalau mau "neraca akuntansi lengkap" (COA + jurnal debit/kredit + hutang/piutang/inventory), perlu modul double-entry.
- Input saldo awal saat ini disimpan di browser (localStorage) per tanggal.

## API (Ringkas)

Auth:

- `POST /login` -> `{ token }`
- `GET /api/me` (auth) -> profile user

Settings:

- `GET /api/settings/store` (auth)
- `PUT /api/settings/store` (admin)

POS/Produk/Order:

- `GET /api/products` (auth)
- `POST /api/orders` (admin/cashier)
- `GET /api/orders?date=YYYY-MM-DD` (admin/finance)
- `GET /api/orders/:id` (admin/finance)

Admin users:

- `GET /api/admin/users` (admin)
- `POST /api/admin/users` (admin)
- `PATCH /api/admin/users/:id` (admin)
- `DELETE /api/admin/users/:id` (admin)

Ledger:

- `GET /api/ledger?date=YYYY-MM-DD` (admin/finance)
- `POST /api/ledger` (admin/finance)
- `DELETE /api/ledger/:id` (admin/finance)

## Development (Opsional)

Saat ini workflow paling simpel adalah pakai Docker Compose.

Kalau mau run frontend di local (hot reload):

```bash
cd frontend
npm install
npm run dev
```

Catatan:

- Pastikan port tidak bentrok dengan container `frontend` (stop `docker compose stop frontend` atau ganti port).
- Frontend default memanggil API ke `http://localhost:8080` (lihat `VITE_API_URL` di `frontend/src/api/api.js`).

## Troubleshooting

- `401 Unauthorized` setelah reset DB:
  - Token di localStorage masih ada tapi user di DB hilang, lakukan logout lalu login lagi.
- Migration baru tidak jalan:
  - Jalankan `./scripts/db_apply_migrations.sh` atau `docker compose down -v` (hapus data).
- Cetak blank / logo tidak muncul:
  - Coba ulang print. Sistem sudah menunggu image load, tapi beberapa driver printer bisa lambat.
- Mau ubah 80mm:
  - Ubah `receiptConfig.printer.receipt` dan `receiptConfig.printer.kitchen` di `frontend/src/config/receipt.js`.
