package main

import (
	"context"
)

// TODO : UPDATE FUNCTION FOR ALL REPOSITORIES

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id int) (*User, error)
	Create(ctx context.Context, user User) (int, error)
	Delete(ctx context.Context, id int) error
	UpdatePassword(ctx context.Context, id int, newHash string) error
	UpdateName(ctx context.Context, id int, newName string) error
}

type TaskRepository interface {
	GetAll(ctx context.Context) ([]Task, error)
	GetById(ctx context.Context, id int) ([]Task, error)
	Create(ctx context.Context, task Task) (int, error)
	Delete(ctx context.Context, id int) error
	UpdateTitle(ctx context.Context, newTitle string, id int) error
	UpdateDescription(ctx context.Context, newDescription string, id int) error
	SwitchTaskStatus(ctx context.Context, id int) error
}
