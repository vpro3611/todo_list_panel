# TODO List — Backend API Service

This repository contains a Go-based backend API for a TODO list application (user + task management). This README focuses on the backend only: how to configure, run, authenticate, use the API, deploy with Docker, and the database schema.

Quick highlights:
- Language: Go 1.25
- HTTP router: chi
- Auth: JWT (HMAC, 24h expiry)
- DB: PostgreSQL (pgxpool)
- Project layout: handlers → services → repositories (Postgres)

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment variables](#environment-variables)
- [Run locally](#run-locally)
- [Docker](#docker)
- [Authentication (JWT)](#authentication-jwt)
- [API Reference](#api-reference)
    - [Public endpoints](#public-endpoints)
    - [Authenticated endpoints: /me](#authenticated-endpoints-me)
    - [Admin endpoints: /admin](#admin-endpoints-admin)
- [Request/Response examples](#requestresponse-examples)
- [Database schema](#database-schema)
- [Migrations](#migrations)
- [Development notes & troubleshooting](#development-notes--troubleshooting)

---

## Prerequisites

- Go 1.25.x
- PostgreSQL (any recent version). Connection via `DATABASE_URL`.
- Docker (optional, for containerized deployment)
- Optional: `psql` or any DB client for running migrations

---

## Environment variables

The app loads environment variables from a file named `config.env` (using `github.com/joho/godotenv`). You may instead set them in your environment or pass via Docker.

Required:
- DATABASE_URL — PostgreSQL connection string. Example: `postgres://user:password@host:5432/dbname?sslmode=disable`
- JWT_SECRET — secret used to sign JWT tokens (HMAC SHA256). Keep this secret safe.

Optional:
- PORT — port the server listens on (defaults to `8080` if not set)

Create `config.env` in the repository root (example):

```
DATABASE_URL=postgres://postgres:pass@localhost:5432/todo_db?sslmode=disable
JWT_SECRET=replace_with_a_strong_secret
PORT=8080
```

> Note: main.go calls `InitConfingEnv()` which uses `godotenv.Load("config.env")`, so filename must be `config.env` or supply env vars in the OS environment.

---

## Run locally

1. Install dependencies and run:

```bash
# from repository root
go mod download
go run .
```

2. Or build and run:

```bash
go build -o todo-backend
./todo-backend
```

On success the server logs "Server started on port: <PORT>" and listens on `:PORT` (default 8080).

---

## Docker

A Dockerfile is provided which builds the Go binary into an Alpine image.

Build the image:

```bash
docker build -t todo-backend .
```

Run the container (recommended approach — map port 8080):

```bash
docker run -d \
  --name todo-backend \
  -e DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=disable" \
  -e JWT_SECRET="replace_with_a_strong_secret" \
  -e PORT=8080 \
  -p 8080:8080 \
  todo-backend
```

Important note about ports:
- The Go app listens on `PORT` (default 8080).
- The included `Dockerfile` contains `EXPOSE 6969` (a metadata instruction). Container port mapping is controlled by `-p`. Use `-p 8080:8080` (or whichever PORT you set) when running the container.

Stop & remove:

```bash
docker stop todo-backend && docker rm todo-backend
```

---

## Authentication (JWT)

- Login endpoint issues a JWT signed with `JWT_SECRET`.
- JWT claims:
    - user_id (int)
    - role (string) — typically `"admin"` or `"user"`
    - registered claims include `exp` (expires 24 hours from issuance)

Include the token in requests to protected endpoints with header:

```
Authorization: Bearer <token>
```

Middleware:
- `JWTmiddleware` verifies token and injects claims into request context.
- `AdminOnly` middleware restricts access to users with role `"admin"`.

Token generation expiration: 24 hours.

---

## API Reference

Base URL: http://localhost:<PORT> (default PORT is 8080)

Models (from code):
- User
    - id: int
    - name: string
    - password: string (hashed on server side)
    - created_at: timestamp
    - updated_at: timestamp
    - role: string ("user" | "admin")
- Task
    - id: int
    - user_id: int
    - title: string
    - description: string
    - is_completed: bool
    - created_at, updated_at: timestamps

All responses are JSON (Content-Type: application/json; charset=UTF-8) unless otherwise noted.

### Public endpoints

- POST /sign-up
    - Description: create a new user
    - Body:
      ```json
      { "name": "alice", "password": "secret123" }
      ```
    - Response: 201 Created (no body)
    - Validations: password must be >= 6 chars

- POST /login
    - Description: authenticate and receive JWT
    - Body:
      ```json
      { "name": "alice", "password": "secret123" }
      ```
    - Response: 200 OK
      ```json
      { "token": "<JWT_TOKEN>" }
      ```

### Authenticated endpoints (/me) — require Authorization header

These endpoints require a valid JWT. The `/me` routes operate on the authenticated user.

- GET /me
    - Returns the current user object (same format as User model).

- PATCH /me/rename
    - Body: { "name": "newname" }
    - Returns: updated user

- PATCH /me/password
    - Body:
      ```json
      { "old_password": "oldPWD", "new_password": "newPWD" }
      ```
    - Returns: updated user (200 OK) on success

- DELETE /me
    - Deletes the current user. Returns JSON containing id and status message.

- GET /me/tasks
    - Returns a list of tasks for the current user.

- POST /me/tasks
    - Body:
      ```json
      { "title": "Buy milk", "description": "2 liters" }
      ```
    - Response: 201 Created and the created task JSON.

- /me/tasks/{id}
    - DELETE -> delete the task (must be owned by the user unless admin)
    - PATCH /title -> body { "title": "New title" } -> returns updated task
    - PATCH /description -> body { "description": "New description" } -> returns updated task
    - PATCH /switch -> toggles task completion -> returns updated task

### Admin endpoints (/admin) — require JWT + AdminOnly

Admin routes allow managing users and all tasks.

- GET /admin/users
    - Returns list of all users.

- POST /admin/users
    - Create a new user (same body as sign-up).
    - Response: 201 Created

- /admin/users/{id}
    - GET -> Get user by id (admin sees all)
    - PATCH /rename -> rename user
    - PATCH /password -> change user password (admin can change others' password)
    - PATCH /role -> toggle role between user/admin. Returns updated user.
    - DELETE -> delete user
    - GET /admin/users/{id}/tasks -> list tasks for specified user
    - POST /admin/users/{id}/tasks -> create task for specified user (body same as create task)

- GET /admin/tasks
    - Returns all tasks.

- /admin/tasks/{id}
    - DELETE -> delete task by id (admin)
    - PATCH /title -> update title (body { "title": "..." })
    - PATCH /description -> update description (body { "description": "..." })
    - PATCH /switch -> toggle is_completed, returns updated task

---

## Request / Response Examples

Replace `<PORT>` with your port and `<JWT_TOKEN>` with the token from login.

1. Sign up

```bash
curl -X POST http://localhost:8080/sign-up \
  -H "Content-Type: application/json" \
  -d '{"name":"alice","password":"secret123"}' -v
# returns HTTP/1.1 201 Created
```

2. Login (get JWT)

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"name":"alice","password":"secret123"}'
# returns: {"token":"<JWT_TOKEN>"}
```

3. Get current user

```bash
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer <JWT_TOKEN>"
# returns user JSON
```

4. Create a task for current user

```bash
curl -X POST http://localhost:8080/me/tasks \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","description":"2 liters"}'
# returns 201 Created and JSON for created task
```

5. Admin: get all users

```bash
curl -X GET http://localhost:8080/admin/users \
  -H "Authorization: Bearer <ADMIN_JWT_TOKEN>"
# returns JSON array of users
```

---

## Database schema

The code expects two main tables: `users` and `tasks`. Below are compatible SQL `CREATE TABLE` statements that match the repository queries and models. Adjust types and extensions as needed for your environment.

NOTE: The schema below includes sensible defaults (timestamps, unique constraint on name, role default).

```sql
-- users table
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'user',
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- tasks table
CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT 'NO DESCRIPTION',
  is_completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
```

The repository code uses queries consistent with these columns:
- `users` columns: id, name, password, created_at, updated_at, role
- `tasks` columns: id, user_id, title, description, is_completed, created_at, updated_at

---

## Migrations

There is a `migrations/` directory (files present in the repo). Use your preferred migration tool (eg. `migrate`, `golang-migrate`, or simple `psql` scripts) to apply the SQL above.

Example (psql):

```bash
# assuming DATABASE_URL is set in env
psql "$DATABASE_URL" -f migrations/20251225143523_create_user_table.up.sql
psql "$DATABASE_URL" -f migrations/20251225143543_create_task_table.up.sql
psql "$DATABASE_URL" -f migrations/20251226183453_add_roles_to_users.up.sql
psql "$DATABASE_URL" -f migrations/20251226210612_make_name_unique.up.sql
```

If you prefer running the SQL directly:

```bash
psql "$DATABASE_URL" <<SQL
-- paste the CREATE TABLE blocks shown above
SQL
```

---

## Development notes & troubleshooting

- config.env: main.go expects a `config.env` file in the project root. Either create it or export env vars globally.
- Database connection: `EstablishDb` requires `DATABASE_URL`; if empty the app will fail with `ErrDBisNotSet`.
- Password hashing / verification: the code encrypts passwords before storing and verifies on login (helpers.go and user_service.go). New passwords must be >= 6 characters.
- Role toggling: `PATCH /admin/users/{id}/role` flips the user's role between `user` and `admin`.
- JWT secret: keep `JWT_SECRET` secret and long enough. Tokens are HMAC-SHA256 signed and valid 24 hours.
- Docker port mismatch: `Dockerfile` contains `EXPOSE 6969` but the server listens on port defined by `PORT` (default 8080). Use `-e PORT=8080 -p 8080:8080` when running the container to avoid confusion.
- If you see errors like "user not found" or permission errors, ensure:
    - You're using a valid JWT token with the correct role in the header.
    - The user / task exists in the DB.
    - For admin routes, the token's role must be `"admin"`.

---


