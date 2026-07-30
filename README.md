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
- PostgreSQL through `pgx/v5`
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
- A database schema compatible with the entities used by this service

Database migrations are not included in this repository yet. See
[doc/database.md](doc/database.md) and [doc/TODO.md](doc/TODO.md) for the current
database implementation notes and remaining work.

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
go run ./cmd/server
```

The API will be available at `http://localhost:3000` unless `PORT` is changed.

Useful checks:

```bash
curl http://localhost:3000/
curl http://localhost:3000/health
```

### Run Tests

```bash
go test ./...
```

## API Overview

Public endpoints:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/` | Welcome response. |
| `GET` | `/health` | Health check with environment and timestamp. |
| `POST` | `/api/auth/register` | Register a user. |
| `POST` | `/api/auth/login` | Log in and receive an access token. |
| `POST` | `/api/auth/refresh` | Refresh the access token from the refresh cookie. |

Protected endpoints require an `Authorization: Bearer <accessToken>` header.
The same API routes are mounted under both `/api` and `/api/v1`.

| Resource | Paths |
| --- | --- |
| Users | `/users`, `/users/:id`, `/users/:id/tasks` |
| Teams | `/teams`, `/teams/:id`, `/teams/:id/members` |
| Projects | `/projects`, `/projects/:id`, `/projects/:id/tasks`, `/projects/:id/teams` |
| Project teams | `/project-teams`, `/project-teams/projects/:projectId/teams` |
| Tasks | `/tasks`, `/tasks/:id`, `/tasks/:id/status` |
| Task assignments | `/task-assignments`, `/task-assignments/tasks/:taskId/assignments` |
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
the main HTTP, service, auth, and PostgreSQL store layers, while versioned
database migrations and broader integration coverage are still tracked as
remaining work in [doc/TODO.md](doc/TODO.md).

## Contributing

Before opening a pull request, run:

```bash
go test ./...
```

Keep changes scoped, include tests for behavior changes, and update this README
or the files in `doc/` when behavior, configuration, or setup steps change.

## License

No license file is currently included in this repository.
