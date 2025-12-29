package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserPgRepository struct {
	pool *pgxpool.Pool
}

func NewUserPgRepository(pool *pgxpool.Pool) *UserPgRepository {
	return &UserPgRepository{
		pool: pool,
	}
}

func (ur *UserPgRepository) GetAll(ctx context.Context) ([]User, error) {
	rows, err := ur.pool.Query(ctx, "SELECT id, name, password, created_at, updated_at, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if users == nil {
		return nil, ErrNoUsers
	}
	return users, nil
}

func (ur *UserPgRepository) GetById(ctx context.Context, id int, actorId int, actorRole string) (*User, error) {
	var u User
	query := "SELECT id, name, password, created_at, updated_at, role FROM users WHERE id = $1 AND (id = $2 OR $3 = 'admin')"
	err := ur.pool.QueryRow(ctx, query, id, actorId, actorRole).Scan(&u.Id,
		&u.Name,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (ur *UserPgRepository) Create(ctx context.Context, user User) (int, error) {
	var id int
	err := ur.pool.QueryRow(ctx, "INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id", user.Name, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ur *UserPgRepository) UpdatePassword(ctx context.Context, id int, newHash string, actorId int, actorRole string) error {
	query := "UPDATE users SET password = $1, updated_at = $2 WHERE id = $3 AND ($3 = $4 OR $5 = 'admin')"
	cmdTag, err := ur.pool.Exec(ctx, query, newHash, time.Now(), id, actorId, actorRole)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ur *UserPgRepository) UpdateName(ctx context.Context, id int, newName string, actorId int, actorRole string) error {
	query := "UPDATE users SET name = $1, updated_at = $2 WHERE id = $3 AND ($3 = $4 OR $5 = 'admin')"

	cmdTag, err := ur.pool.Exec(ctx, query, newName, time.Now(), id, actorId, actorRole)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ur *UserPgRepository) Delete(ctx context.Context, id int, actorId int, actorRole string) error {
	query := "DELETE FROM users WHERE id = $1 AND (id = $2 OR $3 = 'admin')"

	cmdTag, err := ur.pool.Exec(ctx, query, id, actorId, actorRole)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ur *UserPgRepository) UpdateRole(ctx context.Context, id int, newRole string) error {
	cmdTag, err := ur.pool.Exec(ctx, "UPDATE users SET role = $1, updated_at = $2 WHERE id = $3", newRole, time.Now(), id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrSwitchRole
	}
	return nil
}

func (ur *UserPgRepository) Authenticate(ctx context.Context, name string) (*User, error) {
	var user User
	err := ur.pool.QueryRow(ctx, "SELECT id, name, password, created_at, updated_at, role FROM users WHERE name = $1",
		name).Scan(&user.Id,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
