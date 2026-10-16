# Rewrite Roadmap

## Completed

- [x] Membaca `doc/README.md` repo Go dan mengikuti source of truth dari `../project-manager-be`.
- [x] Membaca `../project-manager-be/doc/api_contract.md` dan `../project-manager-be/doc/prd.md`.
- [x] Mengambil dokumentasi Gin terbaru via Context7 untuk route groups, JSON binding, middleware, dan `httptest`.
- [x] Membuat struktur Go service: `cmd/server`, `internal/config`, `internal/auth`, `internal/domain`, `internal/service`, `internal/store/postgres`, `internal/httpapi`.
- [x] Mount route kompatibel `/api` dan `/api/v1`.
- [x] Implement auth JWT access token, refresh token cookie, rotation, dan logout.
- [x] Implement route dasar user, team, project, project-team, task, task-assignment, comment, attachment metadata, notification, `/health`, dan root; upload avatar masih lewat route legacy `/upload`.
- [x] Unit test awal untuk config, service rules, dan Gin route guard.
- [x] Memecah service layer dari satu file besar menjadi file per domain.
- [x] Memecah PostgreSQL store menjadi file per resource dan helper SQL/scanner.
- [x] Memecah HTTP router/handler menjadi file per domain, middleware, response, cookie, dan request helper.
- [x] Memecah service test berdasarkan domain dan menyatukan fixture store di helper test.
- [x] Memecah kontrak store service menjadi interface kecil per domain yang dikomposisi.

## MVP Stabilization

- [x] Perbaiki unit test token auth yang gagal karena token fixed-time sudah expired saat diverifikasi.
- [x] Jalankan ulang `go test ./...` sampai semua package hijau.
- [x] Implementasikan connection supabase
- [x] Tambahkan migration SQL versioned dengan `golang-migrate` berdasarkan target schema shared contract.
- [x] Refactor PostgreSQL store dari `pgx` ke `sqlx` sesuai keputusan stack Go.
- [x] Sediakan route register, login, refresh, logout, `GET /auth/me`, dan compatibility alias `/api`.
- [x] Terapkan response envelope, pagination, dan error shape dasar; parity seluruh payload dan status endpoint masih perlu diaudit.
- [x] Jalankan backend Go dengan PostgreSQL lokal/container, apply migration, lalu validasi smoke test endpoint MVP.
- [x] Tambahkan integration test database untuk repository dan workflow utama.
- [x] Perkuat resource-scoped authorization untuk project, task, comment, attachment, team, dan user/admin flow.
- [x] Stabilkan `PATCH /api/v1/tasks/{task}` untuk edit task dan perpindahan satu task berdasarkan `position`, termasuk otorisasi per field dan test database.
- [x] Simpan refresh token sebagai hash, rotasi saat refresh, dan sembunyikan password hash dari serialisasi `User`.
- [ ] Perkuat auth: rate limit register/login/refresh, rotasi refresh token atomik saat request bersamaan, serta cookie `SameSite`/domain dan CORS yang sesuai untuk frontend/backend lintas domain.
- [x] Samakan `POST /auth/register` dengan kontrak: kembalikan `user`, `accessToken`, dan `refreshToken`, set refresh cookie, dan tolak role istimewa pada registrasi publik.
- [ ] Lengkapi profil user: field `timezone`, `lastLoginAt`, `isActive`, aturan admin-only `password`/`isActive`, dan larang perubahan avatar lewat `PATCH /users/{user}`.
- [ ] Implement `POST /users/{user}/avatar` (Self/Admin) dengan upload image yang divalidasi MIME/size, simpan `avatarStorageKey`, dan kembalikan `avatarUrl`; pertahankan `/upload` sebagai alias legacy bila diperlukan.
- [ ] Lengkapi `GET /users/{user}/tasks`: filter status, priority, projectId, sortBy/order, dan task dari additional assignee, dengan scope manager per resource.
- [ ] Lengkapi query `GET /teams` (`search`, `sortBy`, `order`) serta alias request `user_id` untuk penambahan anggota.
- [ ] Lengkapi model dan request/response project (`status`, `startDate`, `dueDate`, `ownerId`) serta filter `GET /projects` (`status`, `ownerId`).
- [ ] Samakan urutan default `GET /projects/sidebar` dengan kontrak (`name asc` atau `updatedAt desc` sesuai keputusan produk); implementasi sekarang `createdAt desc`.
- [ ] Implement `GET /projects/{project}/summary` dengan count status task, total task, dan jumlah team.
- [ ] Implement route target project-team `POST /projects/{project}/teams` dan `DELETE /projects/{project}/teams/{team}`; route `/project-teams` tetap untuk compatibility.
- [ ] Stabilkan workflow user/profile, team, dan project/sidebar untuk kebutuhan navigasi frontend.
- [ ] Stabilkan workflow task board/list/detail, comment, dan attachment upload/download.
- [ ] Samakan response detail task/comment/attachment dengan core model kontrak, termasuk relasi user/project yang diperlukan frontend, `avatarUrl`, `downloadUrl`, dan alias `size` tanpa mengekspos metadata storage privat.
- [ ] Stabilkan workflow notification dan dashboard berdasarkan kebutuhan MVP frontend.
- [ ] Izinkan ProjectManager+ yang punya akses project menghapus attachment sesuai kontrak; saat ini hanya uploader yang bisa.
- [ ] Implement attachment multipart `file` ke storage lokal/R2 sesuai environment, URL download bertanda tangan/proxy saat request, `GET /attachments/{attachment}/download`, dan hapus object ketika attachment dihapus; batasi MIME dan ukuran file.
- [ ] Lengkapi `GET /notifications` dengan `page`, `limit`, `unreadOnly`, `type`, pagination, `recipientId`, dan `readAt`; set `readAt` saat notifikasi dibaca serta hitung `meta.unreadCount` dari seluruh notifikasi user, bukan hanya halaman aktif.
- [ ] Lengkapi `GET /api/v1/health` publik dan payload diagnostik `/health` (`database`, `queue`, `storage`; metadata deployment bila tersedia) tanpa data rahasia.
- [ ] Siapkan konfigurasi runtime kontrak untuk `CORS_ORIGINS`, `COOKIE_DOMAIN`, `STORAGE_DRIVER`, R2, dan signed URL TTL; jangan gunakan secret JWT default di production.
- [ ] Audit parity response terhadap shared API contract/Postman setelah frontend dicoba.
- [ ] Tambahkan dokumentasi OpenAPI/Swagger untuk menggantikan endpoint `/api-docs` Node.

## Dashboard Feature

- [x] Implement `GET /api/v1/dashboard` dengan scope authenticated user's visible projects/tasks.
- [x] Pastikan dashboard ikut tersedia lewat compatibility alias `/api/dashboard`.
- [x] Kembalikan response dashboard sebagai object langsung sesuai frontend `dashboardService`, bukan envelope `{ data: ... }`, kecuali frontend ikut diubah untuk unwrap envelope.
- [x] Kembalikan payload contract: `totalActiveProjects`, `taskCountPerStatus`, `activeProgress`, `inReview`, `dueSoon`, `overdue`, `overdueTaskCount`, dan `workloadPerMember`.
- [x] Hitung `activeProgress.doing` dari `in_progress`, `activeProgress.todo` dari `todo`, `total = todo + doing`, `ratio` 4 desimal, dan `percentage` 2 desimal.
- [x] Hitung `dueSoon` untuk task belum done yang jatuh tempo hari ini sampai 7 hari ke depan.
- [x] Hitung `overdue` untuk task belum done yang jatuh tempo sebelum hari ini, dan samakan `overdueTaskCount` sebagai compatibility alias.
- [x] Kembalikan list max 5 item untuk `recentlyUpdatedTasks`, `recentTasks`, `upcomingDeadlines`, dan `highPriorityTasks`.
- [x] Kembalikan `latestUpdates` max 5 activity log item, atau `[]` sampai activity log P1 tersedia.
- [x] Terapkan resource-scoped authorization pada semua query dashboard.
- [x] Tambahkan repository/service test untuk agregasi dashboard.
- [x] Tambahkan integration test database untuk workflow dashboard utama.

## P1

- [x] Catat activity dasar create/update/delete task dan sediakan `GET /tasks/{task}/activity` dengan scope task.
- [ ] Lengkapi activity log untuk aksi project, assignment, comment, dan perubahan penting lain; sediakan `GET /activity-logs` dan `GET /projects/{project}/activity` sesuai scope kontrak, lalu tampilkan itemnya di `dashboard.latestUpdates`.
- [x] Implement `POST /api/v1/tasks/reorder` untuk reorder batch atomik sesuai API contract.
- [ ] Implement unified assignment.
- [ ] Implement mention berbasis user ID yang stabil; parser `@name` saat ini hanya memakai nama satu kata dan dapat ambigu.
- [ ] Implement checklist.
- [ ] Implement `POST /exports/project-report` (CSV project report/task list) dan `GET /exports/{export}` dengan scope pemilik/project manager.
- [ ] Implement CRUD `/webhook-endpoints`, `GET /webhook-deliveries`, delivery via queue, HMAC signature, dan retry kegagalan.

## P2

- [ ] Implement realtime notification.
- [ ] Implement insight lanjutan.

## Sudah ada, tetapi belum tercatat eksplisit

- [x] `GET /projects/sidebar` mengembalikan project visible beserta `openTaskCount`.
- [x] `GET /tasks` mendukung mode board tanpa pagination, mode list dengan pagination, alias `project_id`, dan filter dasar status/priority/assignee/due date.
- [x] Route assignment task target (`POST /tasks/{task}/assignments`, `DELETE /tasks/{task}/assignments/{user}`) sudah tersedia bersama route legacy `/task-assignments`.
- [x] Comment `@name` dan assignment task sudah membuat notifikasi; `GET /notifications` saat ini mengembalikan `meta.unreadCount` untuk seluruh item yang diambil.
