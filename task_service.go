package main

import (
	"context"
	"errors"
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
		return nil, errors.New("id must be greater than 0")
	}
	tasks, err := ts.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		return nil, errors.New("no tasks found")
	}

	return tasks, err
}

func (ts *TaskService) CreateNewTask(ctx context.Context, userId int, title string, description string) (int, error) {
	if userId < 1 {
		return 0, errors.New("user id must be greater than 0")
	}
	if title == "" {
		return 0, errors.New("title must be not empty")
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
			return 0, errors.New("user with this id does not exist")
		}
		return 0, err
	}
	return id, nil
}

func (ts *TaskService) DeleteTask(ctx context.Context, id int) error {
	if id < 1 {
		return errors.New("id must be greater than 0")
	}
	return ts.repo.Delete(ctx, id)
}

func (ts *TaskService) UpdateTitle(ctx context.Context, newTitle string, id int) error {
	if id < 1 {
		return errors.New("id must be greater than 0")
	}
	if newTitle == "" {
		return errors.New("title must be not empty")
	}

	return ts.repo.UpdateTitle(ctx, newTitle, id)
}

func (ts *TaskService) UpdateDescription(ctx context.Context, newDescription string, id int) error {
	if id < 1 {
		return errors.New("id must be greater than 0")
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
		return errors.New("id must be greater than 0")
	}

	return ts.repo.SwitchTaskStatus(ctx, id)
}
