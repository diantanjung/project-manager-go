# Database Implementation - Golang Backend

> Dokumen ini hanya berlaku untuk implementasi `project-manager-go`.
> Product behavior, endpoint, payload, response, enum, dan authorization tetap mengikuti shared docs di `../project-manager-be/doc`.

## Stack Decision

- Database utama: PostgreSQL.
- Database access: `sqlx` di atas `database/sql`.
- Migration: SQL versioned files dengan `golang-migrate`.
- ORM seperti GORM/Ent tidak dipakai sebagai abstraction utama.

## Rationale

- SQL tetap eksplisit dan mudah diaudit.
- `sqlx` mengurangi boilerplate scan struct tanpa menyembunyikan query.
- `golang-migrate` menjaga schema change sebagai artifact versioned yang bisa direview.
- API contract tetap tidak bergantung pada pilihan library internal.

## Implementation Notes

- Simpan migration di folder khusus repo Go, misalnya `migrations/`.
- Semua query runtime wajib memakai context-aware method.
- Gunakan parameterized query, bukan string concatenation untuk input user.
- Transaction boundary ditentukan di service/repository sesuai operasi bisnis.
- Response API tetap `camelCase`; kolom database boleh `snake_case`.
