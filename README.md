# Project Manager Go

Project Manager Go is a Go rewrite of the Project Manager backend API. It
provides authentication, role-based access control, project and team management,
task tracking, comments, attachments, notifications, and PostgreSQL-backed
persistence.

This repository is intended to host the Go backend implementation. Shared
product behavior, API contract, and PRD references currently live in the
original backend project described in [doc/README.md](doc/README.md).

## Tech Stack

- Go 1.26.5
- Gin for HTTP routing and middleware
- PostgreSQL through `sqlx` with the `pgx` database/sql driver
- JWT access and refresh tokens
- Cookie-based refresh token rotation
- `golang.org/x/crypto/bcrypt` password hashing

## Features

- User registration, login, refresh, and logout
- JWT-protected REST API mounted at both `/api` and `/api/v1`
- Role-based authorization for admin, product owner, project manager, and team
  member workflows
- CRUD endpoints for users, teams, projects, and tasks
- Project-team assignment and task assignment endpoints
- Comments, attachments, avatar upload, and notification endpoints
- Pagination and filtering for list endpoints
- Health check endpoint at `/health`
- Graceful HTTP server shutdown on `SIGINT` and `SIGTERM`

## Getting Started

### Prerequisites

- Go 1.26.5 or newer
- PostgreSQL
- A database reachable through `DATABASE_URL`

Versioned SQL migrations live in [migrations](migrations). The database helper
command uses `sqlx` for connection checks and `golang-migrate` for schema
changes.

### Configuration

The server reads environment variables directly and also loads a local `.env`
file when present.

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `DATABASE_URL` | Yes | - | PostgreSQL connection URL. |
| `PORT` | No | `3000` | HTTP server port. |
| `NODE_ENV` | No | `development` | Set to `production` to enable Gin release mode. |
| `FRONTEND_URL` | No | `http://localhost:5173` | Allowed CORS origin. |
| `JWT_SECRET` | No | `supersecret` | Access token signing secret. Override outside local development. |
| `JWT_EXPIRES_IN` | No | `15m` | Access token TTL. Supports `s`, `m`, `h`, and `d`. |
| `JWT_REFRESH_SECRET` | No | `superrefreshsecret` | Refresh token signing secret. Override outside local development. |
| `JWT_REFRESH_EXPIRES_IN` | No | `7d` | Refresh token TTL. Supports `s`, `m`, `h`, and `d`. |
| `UPLOAD_DIR` | No | `uploads` | Directory served from `/uploads`. |

Example local `.env`:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/project_manager?sslmode=disable
PORT=3000
FRONTEND_URL=http://localhost:5173
JWT_SECRET=replace-me
JWT_REFRESH_SECRET=replace-me-too
```

### Run Locally

```bash
go mod download
go run ./cmd/db ping
go run ./cmd/db migrate-up
go run ./cmd/db seed
go run ./cmd/server
```

The API will be available at `http://localhost:3000` unless `PORT` is changed.

Useful checks:

```bash
curl http://localhost:3000/
curl http://localhost:3000/health
```

Migration commands:

```bash
go run ./cmd/db migrate-version
go run ./cmd/db migrate-down
go run ./cmd/db migrate-steps 1
go run ./cmd/db seed
```

Set `MIGRATIONS_SOURCE` when running from outside the repository root. The
default is `file://migrations`.

### Production environment

Use [`.env.production.example`](.env.production.example) as the variable list.
The local `.env.production` file is ignored by Git. Set `DATABASE_URL` to the
production PostgreSQL URL with TLS enabled and `FRONTEND_URL` to the exact HTTPS
origin of the frontend (no trailing slash). Use two different, random JWT
secrets. Keep `UPLOAD_DIR` on persistent storage if uploads must survive a
redeploy. The refresh cookie uses `Secure` mode when `NODE_ENV=production`, so
serve the API over HTTPS. Its `SameSite=Strict` setting means browsers send the
cookie only for same-site requests. For cross-site deployments such as Vercel
frontend and Render API, login and refresh responses also return `refreshToken`.
The frontend stores that fallback in `sessionStorage` and sends it in the JSON
body of refresh and logout requests. This keeps sessions after page reloads even
when cross-site cookies are blocked. Because JavaScript can read the fallback
token, protect the frontend against XSS. The HttpOnly cookie remains the preferred
source when the browser sends it.

The app automatically reads `.env`, not `.env.production`. On a server, copy the
completed production file to `.env` in the process working directory, or enter
the same variables in the hosting platform's environment or secret settings.
Run `go run ./cmd/db migrate-up` with those variables before starting
`go run ./cmd/server`. Do not run `go run ./cmd/db seed` in production; it loads
demo data. Set `MIGRATIONS_SOURCE` if migrations are outside the working
directory.

### Run Tests

```bash
go test ./...
```

Database integration tests are gated behind the `integration` build tag and use
`TEST_DATABASE_URL` when set, falling back to `DATABASE_URL` from `.env`. They
run migrations first and roll test data back with transactions.

```bash
go test -tags=integration ./internal/store/postgres -run Integration -v
```

## API Overview

Public endpoints:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/` | Welcome response. |
| `GET` | `/health` | Health check with environment and timestamp. |
| `POST` | `/api/auth/register` | Register a user. |
| `POST` | `/api/auth/login` | Log in and receive access and refresh tokens; also set the refresh cookie. |
| `POST` | `/api/auth/refresh` | Rotate the refresh token from the cookie or JSON `refreshToken` body. |
| `GET` | `/api/auth/me` | Return the current authenticated user. |
| `POST` | `/api/auth/logout` | Revoke the refresh token from the cookie or JSON body and clear the cookie. |

Protected endpoints require an `Authorization: Bearer <accessToken>` header.
The same API routes are mounted under both `/api` and `/api/v1`.

| Resource | Paths |
| --- | --- |
| Users | `/users`, `/users/:id`, `/users/:id/tasks` |
| Teams | `/teams`, `/teams/:id`, `/teams/:id/members` |
| Projects | `/projects`, `/projects/:id`, `/projects/:id/tasks`, `/projects/:id/teams` |
| Project teams | `/project-teams`, `/project-teams/projects/:projectId/teams` |
| Tasks | `/tasks` (list or board), `/tasks/reorder`, `/tasks/:id`, `/tasks/:id/status`, `/tasks/:id/activity` |
| Task assignments | `/tasks/:id/assignments`, `/tasks/:id/assignments/:userId`; legacy `/task-assignments` routes |
| Comments | `/tasks/:id/comments`, `/comments/:id` |
| Attachments | `/tasks/:id/attachments`, `/attachments/:id` |
| Uploads | `/upload`, `/uploads/*` |
| Notifications | `/notifications`, `/notifications/:id/read`, `/notifications/read-all` |

## Roles

The domain model defines these application roles:

- `admin`
- `productOwner`
- `projectManager`
- `teamMember`

Authorization rules are enforced in HTTP middleware. For example, creating users
requires `admin`, creating projects requires `projectManager`, and deleting
projects requires `productOwner`.

## Project Layout

```text
cmd/db/                  Database ping and migration commands
cmd/server/              HTTP server entry point
internal/auth/           JWT token creation and validation
internal/config/         Environment and .env loading
internal/domain/         Domain models, enums, and application errors
internal/httpapi/        Gin router, middleware, handlers, and responses
internal/service/        Business logic and authorization-aware workflows
internal/store/postgres/ PostgreSQL persistence layer
doc/                     Implementation notes and rewrite tracking
```

## Development Status

This is an active rewrite of an existing backend. The Go service already includes
the main HTTP, service, auth, PostgreSQL store layers, and initial versioned SQL
migrations, while broader integration coverage and future milestones are tracked
in [doc/ROADMAP.md](doc/ROADMAP.md).

## Contributing

Before opening a pull request, run:

```bash
go test ./...
```

Keep changes scoped, include tests for behavior changes, and update this README
or the files in `doc/` when behavior, configuration, or setup steps change.

## License

No license file is currently included in this repository.
