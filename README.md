# POS Warkop

POS kasir sederhana untuk warkop/cafe: POS + cetak struk & kitchen ticket + dashboard finance + setting toko + manajemen user & produk.

## Fitur Unggulan Baru 🚀

- **Manajemen Produk & Gambar**: Upload foto menu langsung dari browser.
- **Integrasi Midtrans (QRIS Dinamis)**: Pembayaran otomatis via Midtrans Snap (QRIS, VA, dsb). Tidak perlu input nominal manual saat scan.
- **Webhook Status**: Status pesanan otomatis berubah jadi 'Lunas' saat pembayaran Midtrans berhasil.

## Fitur Utama

- **POS (Kasir)**
  - Cari menu, tambah ke keranjang, atur qty.
  - Pembayaran Cash / QRIS (Manual) / **Midtrans (Otomatis)**.
  - Dine In / Take Away, input meja, jumlah tamu, nama pemesan.
  - Cetak Struk customer & Ticket kitchen.
- **Setting (Admin)**
  - Profil toko & Upload logo.
  - **Manajemen Produk**: Tambah/Edit/Hapus menu + foto.
  - Manajemen user: cashier, finance, admin (superadmin).
- **Finance (Admin/Finance)**
  - Trace transaksi harian + detail order.
  - Pembukuan manual (income/expense).
  - Neraca harian sederhana.

## Roles (Akses)

- `cashier`: hanya POS, bisa buat order.
- `finance`: hanya Finance, akses laporan dan pembukuan.
- `admin`: akses semua fitur (POS + Finance + Setting + Users + Products).

## Tech Stack

- Backend: Go 1.22, Gorilla Mux, Postgres, JWT.
- Frontend: Vue 3, Vite, TailwindCSS, Pinia.
- Payment: Midtrans Snap SDK integration.

## Environment Variables

Copy `.env.example` menjadi `.env` dan sesuaikan nilainya.

- `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB`: Kredensial DB.
- `JWT_SECRET`: Secret untuk token keamanan.
- `MIDTRANS_SERVER_KEY`: Server key dari Dashboard Midtrans.
- `MIDTRANS_CLIENT_KEY`: Client key dari Dashboard Midtrans.

## Quick Start (Docker)

```bash
docker compose up -d --build
```

Buka:
- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- Default Admin: `admin@pos.com` / `123456`

## Database & Migrations

File migration baru:
- `006_product_images.sql`: Dukungan upload foto produk.
- `007_midtrans_support.sql`: Kolom status & token pembayaran Midtrans.

Gunakan script berikut jika DB sudah ada datanya:
```bash
./scripts/db_apply_migrations.sh
```

## Receipt & Kitchen Printing

Konfigurasi di `frontend/src/config/receipt.js`. Mendukung printer thermal 58mm (default) dan 80mm.

## API Endpoints (Baru)

- `POST /api/webhooks/midtrans`: Webhook notifikasi pembayaran.
- `POST /api/admin/products`: Tambah produk + foto.
- `PUT /api/admin/products/:id`: Update produk + foto.
- `DELETE /api/admin/products/:id`: Hapus produk.

---
**Happy Brewing! ☕**
