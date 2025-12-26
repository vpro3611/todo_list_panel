package main

import (
	"context"
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

func (ts *TaskService) GetTaskById(ctx context.Context, id int) ([]Task, error) {
	if id < 1 {
		return nil, ErrIdMustBeGtZero
	}
	tasks, err := ts.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		return nil, ErrNoTasks
	}

	return tasks, err
}

func (ts *TaskService) CreateNewTask(ctx context.Context, userId int, title string, description string) (int, error) {
	if userId < 1 {
		return 0, ErrIdMustBeGtZero
	}
	if title == "" {
		return 0, ErrEmptyTitle
	}

	var desc string

	if description == "" {
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

func (ts *TaskService) DeleteTask(ctx context.Context, id int) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	return ts.repo.Delete(ctx, id)
}

func (ts *TaskService) UpdateTitle(ctx context.Context, newTitle string, id int) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	if newTitle == "" {
		return ErrEmptyTitle
	}

	return ts.repo.UpdateTitle(ctx, newTitle, id)
}

func (ts *TaskService) UpdateDescription(ctx context.Context, newDescription string, id int) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	var desc string

	if newDescription == "" {
		desc = "NO DESCRIPTION"
	} else {
		desc = newDescription
	}

	return ts.repo.UpdateDescription(ctx, desc, id)
}

func (ts *TaskService) SwitchTaskStatus(ctx context.Context, id int) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}

	return ts.repo.SwitchTaskStatus(ctx, id)
}
