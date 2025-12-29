package main

import (
	"context"
)

// TODO : UPDATE FUNCTION FOR ALL REPOSITORIES

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id int, actorId int, actorRole string) (*User, error)
	Create(ctx context.Context, user User) (int, error)
	Delete(ctx context.Context, id int, actorId int, actorRole string) error
	UpdatePassword(ctx context.Context, id int, newHash string, actorId int, actorRole string) error
	UpdateName(ctx context.Context, id int, newName string, actorId int, actorRole string) error
	UpdateRole(ctx context.Context, id int, newRole string) error
	Authenticate(ctx context.Context, name string) (*User, error)
}

type TaskRepository interface {
	GetAll(ctx context.Context) ([]Task, error)
	GetByUserId(ctx context.Context, id int, actorId int, actorRole string) ([]Task, error)
	Create(ctx context.Context, task Task) (int, error)
	Delete(ctx context.Context, id int, actorId int, actorRole string) error
	UpdateTitle(ctx context.Context, newTitle string, id int, actorId int, actorRole string) error
	UpdateDescription(ctx context.Context, newDescription string, id int, actorId int, actorRole string) error
	SwitchTaskStatus(ctx context.Context, id int, actorId int, actorRole string) error
}
