package main

import (
	"context"
	"strings"
)

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (ts *TaskService) GetAllTasks(ctx context.Context) ([]Task, error) {
	return ts.repo.GetAll(ctx)
}

func (ts *TaskService) GetTaskById(ctx context.Context, id int, actorId int, actorRole string) ([]Task, error) {
	if id < 1 {
		return nil, ErrIdMustBeGtZero
	}
	return ts.repo.GetByUserId(ctx, id, actorId, actorRole)
}

func (ts *TaskService) CreateNewTask(ctx context.Context, userId int, title string, description string) (int, error) {
	if userId < 1 {
		return 0, ErrIdMustBeGtZero
	}
	if strings.TrimSpace(title) == "" {
		return 0, ErrEmptyTitle
	}

	var desc string

	if strings.TrimSpace(description) == "" {
		desc = "NO DESCRIPTION"
	} else {
		desc = description
	}

	newTask := Task{
		UserId:      userId,
		Title:       title,
		Description: desc,
	}

	id, err := ts.repo.Create(ctx, newTask)
	if err != nil {
		if IsForeignKeyViolation(err) {
			return 0, ErrNoUserWithThisId
		}
		return 0, err
	}
	return id, nil
}

func (ts *TaskService) DeleteTask(ctx context.Context, id int, actorId int, actorRole string) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	return ts.repo.Delete(ctx, id, actorId, actorRole)
}

func (ts *TaskService) UpdateTitle(ctx context.Context, newTitle string, id int, actorId int, actorRole string) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	if newTitle == "" {
		return ErrEmptyTitle
	}

	return ts.repo.UpdateTitle(ctx, newTitle, id, actorId, actorRole)
}

func (ts *TaskService) UpdateDescription(ctx context.Context, newDescription string, id int, actorId int, actorRole string) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	var desc string

	if newDescription == "" {
		desc = "NO DESCRIPTION"
	} else {
		desc = newDescription
	}

	return ts.repo.UpdateDescription(ctx, desc, id, actorId, actorRole)
}

func (ts *TaskService) SwitchTaskStatus(ctx context.Context, id int, actorId int, actorRole string) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}

	return ts.repo.SwitchTaskStatus(ctx, id, actorId, actorRole)
}

func (ts *TaskService) GetTaskByItsId(ctx context.Context, taskId int, actorId int, actorRole string) (*Task, error) {
	if taskId < 1 {
		return nil, ErrIdMustBeGtZero
	}
	return ts.repo.GetTaskById(ctx, taskId, actorId, actorRole)
}
