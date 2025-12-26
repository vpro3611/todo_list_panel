package main

import (
	"context"
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

func (ur *UserPgRepository) GetById(ctx context.Context, id int) (*User, error) {
	var u User
	err := ur.pool.QueryRow(ctx, "SELECT id, name, password, created_at, updated_at, role FROM users WHERE id = $1", id).Scan(&u.Id,
		&u.Name,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Role)
	if err != nil {
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

func (ur *UserPgRepository) UpdatePassword(ctx context.Context, id int, newHash string) error {
	cmdTag, err := ur.pool.Exec(ctx, "UPDATE users SET password = $1, updated_at = $2 WHERE id = $3", newHash, time.Now(), id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ur *UserPgRepository) UpdateName(ctx context.Context, id int, newName string) error {
	cmdTag, err := ur.pool.Exec(ctx, "UPDATE users SET name = $1, updated_at = $2 WHERE id = $3", newName, time.Now(), id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ur *UserPgRepository) Delete(ctx context.Context, id int) error {
	_, err := ur.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
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
