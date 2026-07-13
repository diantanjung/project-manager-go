# Rewrite TODO

## Completed

- [x] Membaca `doc/README.md` repo Go dan mengikuti source of truth dari `../project-manager-be`.
- [x] Membaca `../project-manager-be/doc/api_contract.md` dan `../project-manager-be/doc/prd.md`.
- [x] Mengambil dokumentasi Gin terbaru via Context7 untuk route groups, JSON binding, middleware, dan `httptest`.
- [x] Membuat struktur Go service: `cmd/server`, `internal/config`, `internal/auth`, `internal/domain`, `internal/service`, `internal/store/postgres`, `internal/httpapi`.
- [x] Mount route kompatibel `/api` dan `/api/v1`.
- [x] Implement auth JWT access token, refresh token cookie, rotation, dan logout.
- [x] Implement user, team, project, project-team, task, task-assignment, comment, attachment, upload avatar, notification, health, dan root routes.
- [x] Unit test awal untuk config, token auth, service rules, dan Gin route guard.
- [x] Memecah service layer dari satu file besar menjadi file per domain.
- [x] Memecah PostgreSQL store menjadi file per resource dan helper SQL/scanner.
- [x] Memecah HTTP router/handler menjadi file per domain, middleware, response, cookie, dan request helper.
- [x] Memecah service test berdasarkan domain dan menyatukan fixture store di helper test.
- [x] Memecah kontrak store service menjadi interface kecil per domain yang dikomposisi.

## Remaining

- [ ] Jalankan aplikasi melawan database PostgreSQL nyata dengan schema/migration dari backend Node.
- [ ] Tambahkan integration test database dengan test container atau database test dedicated.
- [ ] Tambahkan Swagger/OpenAPI generation untuk menggantikan endpoint `/api-docs` Node.
- [ ] Tambahkan rate limiting login dan hardening upload MIME sniffing.
- [ ] Audit parity response detail terhadap Postman collection setelah frontend dicoba.
