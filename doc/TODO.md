# TODO

## Completed

- [x] Membaca `doc/README.md` repo Go dan mengikuti source of truth dari `../project-manager-be`.
- [x] Membaca `../project-manager-be/doc/api_contract.md` dan `../project-manager-be/doc/prd.md`.
- [x] Mengambil dokumentasi Gin terbaru via Context7 untuk route groups, JSON binding, middleware, dan `httptest`.
- [x] Membuat struktur Go service: `cmd/server`, `internal/config`, `internal/auth`, `internal/domain`, `internal/service`, `internal/store/postgres`, `internal/httpapi`.
- [x] Mount route kompatibel `/api` dan `/api/v1`.
- [x] Implement auth JWT access token, refresh token cookie, rotation, dan logout.
- [x] Implement user, team, project, project-team, task, task-assignment, comment, attachment, upload avatar, notification, health, dan root routes.
- [x] Unit test awal untuk config, service rules, dan Gin route guard.
- [x] Memecah service layer dari satu file besar menjadi file per domain.
- [x] Memecah PostgreSQL store menjadi file per resource dan helper SQL/scanner.
- [x] Memecah HTTP router/handler menjadi file per domain, middleware, response, cookie, dan request helper.
- [x] Memecah service test berdasarkan domain dan menyatukan fixture store di helper test.
- [x] Memecah kontrak store service menjadi interface kecil per domain yang dikomposisi.

## Remaining

- [x] Perbaiki unit test token auth yang gagal karena token fixed-time sudah expired saat diverifikasi.
- [x] Jalankan ulang `go test ./...` sampai semua package hijau.
- [x] Implementasikan connection supabase
- [x] Tambahkan migration SQL versioned dengan `golang-migrate` berdasarkan target schema shared contract.
- [x] Refactor PostgreSQL store dari `pgx` ke `sqlx` sesuai keputusan stack Go.
- [ ] Pastikan auth/session MVP lengkap: register, login, refresh, logout, `GET /auth/me`, dan compatibility alias `/api`.
- [ ] Rapikan response envelope, pagination, dan error shape agar konsisten dengan shared API contract.
- [ ] Jalankan backend Go dengan PostgreSQL lokal/container, apply migration, lalu validasi smoke test endpoint MVP.
- [ ] Tambahkan integration test database untuk repository dan workflow utama.
- [ ] Perkuat resource-scoped authorization untuk project, task, comment, attachment, team, dan user/admin flow.
- [ ] Perkuat security hardening: safe user response, refresh token hash/rotation, rate limit auth, dan upload MIME/size validation.
- [ ] Stabilkan workflow user/profile, team, dan project/sidebar untuk kebutuhan navigasi frontend.
- [ ] Stabilkan workflow task board/list/detail, comment, dan attachment upload/download.
- [ ] Stabilkan workflow notification dan dashboard berdasarkan kebutuhan MVP frontend.
- [ ] Audit parity response terhadap shared API contract/Postman setelah frontend dicoba.
- [ ] Tambahkan dokumentasi OpenAPI/Swagger untuk menggantikan endpoint `/api-docs` Node.
- [ ] Implement P1 setelah MVP stabil: activity log, reorder atomik, unified assignment, dan checklist.
- [ ] Implement P2 setelah P1 stabil: export, webhook, realtime notification, dan insight lanjutan.
