# POS Warkop

POS kasir sederhana untuk warkop/cafe: POS + KDS (Kitchen Display System) + cetak struk & kitchen ticket + dashboard finance + setting toko + manajemen user & produk.

## Fitur Unggulan Baru 🚀

- **KDS (Kitchen Display System)**: Dashboard khusus dapur untuk memantau dan memproses pesanan secara real-time.
- **Monitoring Dapur (POS Sidebar)**: Kasir dapat memantau status pesanan (pending, cooking, ready, done) langsung dari sidebar POS tanpa pindah halaman.
- **Dashboard Khusus Produk**: Manajemen menu lebih rapi dan lapang di dashboard khusus.
- **Tersedia/Habis (Stok)**: Sembunyikan produk dari POS secara instan jika stok sedang kosong.
- **Integrasi Midtrans (QRIS Dinamis)**: Pembayaran otomatis via Midtrans Snap (QRIS, VA, dsb). Tidak perlu input nominal manual saat scan.
- **Pencarian Transaksi**: Cari order berdasarkan ID atau Nomor Antrian di halaman Finance.
- **Webhook Status**: Status pesanan otomatis berubah jadi 'Lunas' saat pembayaran Midtrans berhasil.

## Fitur Utama

- **POS (Kasir)**
  - Cari menu, tambah ke keranjang, atur qty.
  - **Monitoring Dapur**: Lihat status progres makanan (diantre/dimasak/siap/selesai) secara live via sidebar.
  - Pembayaran Cash / QRIS (Manual) / **Midtrans (Otomatis)**.
  - Dine In / Take Away, input meja, jumlah tamu, nama pemesan.
  - Cetak Struk customer & Ticket kitchen.
- **KDS (Kitchen Display System)**
  - Dashboard khusus untuk kru dapur.
  - Update status pesanan (Pending -> Cooking -> Ready -> Done) dengan sekali klik.
  - Auto-refresh untuk pesanan baru.
- **Produk (Dashboard Khusus)**
  - Tambah/Edit/Hapus menu + foto.
  - **Status Stok**: Aktifkan/Matikan menu (Disable menu jika kosong).
- **Setting (Admin)**
  - Profil toko & Upload logo.
  - Manajemen user: cashier, kitchen, finance, admin (superadmin).
- **Finance (Admin/Finance)**
  - **Pencarian Order**: Cari transaksi berdasarkan ID / No Order.
  - Trace transaksi harian + detail order.
  - Pembukuan manual (income/expense).
  - Neraca harian sederhana.

## Roles (Akses)

- `cashier`: akses POS (buat order & monitoring dapur).
- `kitchen`: akses KDS (Kitchen Display System) saja.
- `finance`: akses menu Finance (laporan & pembukuan).
- `admin`: akses semua fitur (POS + KDS + Finance + Setting + Users + Products).

## Tech Stack

- Backend: Go 1.22, Gorilla Mux, Postgres, JWT.
- Frontend: Vue 3 (Composition API), Vite, TailwindCSS, Pinia.
- Payment: Midtrans Snap SDK integration.

## Environment Variables

Copy `.env.example` menjadi `.env` dan sesuaikan nilainya.

- `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB`: Kredensial DB.
- `JWT_SECRET`: Secret untuk token keamanan.
- `MIDTRANS_SERVER_KEY`: Server key dari Dashboard Midtrans.
- `MIDTRANS_CLIENT_KEY`: Client key dari Dashboard Midtrans.

## Cara Menjalankan (Docker - Decoupled)

Sistem sekarang menggunakan arsitektur terpisah untuk keamanan data:

1. **Jalankan Database (Infra)**
   ```bash
   docker compose -f docker-compose.infra.yml up -d
   ```
   *Biarkan ini menyala terus agar data tetap aman.*

2. **Jalankan Aplikasi (Backend & Frontend)**
   ```bash
   docker compose -f docker-compose.app.yml up --build -d
   ```
   *Gunakan ini jika ada perubahan kode atau ingin restart aplikasi.*

3. **Matikan Aplikasi (Tanpa matikan DB)**
   ```bash
   docker compose -f docker-compose.app.yml down
   ```

*Catatan: Gunakan `docker-compose.yml` (default) jika ingin menjalankan semuanya dalam satu lifecycle (resiko data hilang jika `down -v`).*

## Database & Migrations

Aplikasi menggunakan **Automated Migrations** di dalam kode Go. 
- Tidak perlu lagi menjalankan script migrasi manual.
- Backend akan otomatis mengecek dan melakukan migrasi tabel saat pertama kali dijalankan.
- Pastikan folder `migrations/` ada di root project.

## Receipt & Kitchen Printing

Konfigurasi di `frontend/src/config/receipt.js`. Mendukung printer thermal 58mm (default) dan 80mm.

## API Endpoints (Baru/Penting)

- `PATCH /api/orders/:id/status`: Update status dapur (oleh kitchen role).
- `POST /api/webhooks/midtrans`: Webhook notifikasi pembayaran.
- `POST /api/admin/products`: Tambah produk + foto.
- `PUT /api/admin/products/:id`: Update produk + foto.
- `DELETE /api/admin/products/:id`: Hapus produk.

---
**Happy Brewing! ☕**
