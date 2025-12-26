# Awesome Task Manager API

A Go-based RESTful API for managing users and their associated tasks. This project demonstrates a clean architectural approach using services, repositories, and PostgreSQL for persistence.

## Features

- **User Management**: Create, update, and delete users.
- **Task Management**: Create tasks for specific users, update titles/descriptions, and toggle completion status.
- **Graceful Shutdown**: Handles OS signals to shut down the server and database connections safely.
- **Environment Configuration**: Uses `.env` files for database connection settings.

## Tech Stack

- **Language**: Go 1.25+
- **HTTP Router**: [go-chi/chi](https://github.com/go-chi/chi)
- **Database**: PostgreSQL
- **Environment Management**: [joho/godotenv](https://github.com/joho/godotenv)

## Project Structure

- `main.go`: Entry point, handles initialization and graceful shutdown.
- `server.go`: HTTP server setup, routing, and task handlers.
- `task_service.go` / `user_service.go`: Business logic layer.
- `task_repo_pg.go` / `user_repo_pg.go`: PostgreSQL implementation of the repositories.
- `models.go`: Core data structures (User, Task).
- `interfaces.go`: Definition of repository interfaces for decoupling.
- `migrations/`: Database schema and migration files.

## Getting Started

### Prerequisites

- Go installed on your machine.
- PostgreSQL database running.

### Configuration

Create a `database.env` file in the root directory with the following variables:
```env
DB_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
```