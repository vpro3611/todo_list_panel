package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TaskPgRepository struct {
	pool *pgxpool.Pool
}

func NewTaskPgRepository(pool *pgxpool.Pool) *TaskPgRepository {
	return &TaskPgRepository{
		pool: pool,
	}
}

func (tr *TaskPgRepository) GetAll(ctx context.Context) ([]Task, error) {
	var tasks []Task

	rows, err := tr.pool.Query(ctx, "SELECT id, user_id, title, description, is_completed, created_at, updated_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id,
			&t.UserId,
			&t.Title,
			&t.Description,
			&t.IsCompleted,
			&t.CreatedAt,
			&t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if tasks == nil {
		return nil, ErrNoTasks
	}

	return tasks, nil
}

func (tr *TaskPgRepository) GetByUserId(ctx context.Context, id int) ([]Task, error) {
	var tasks []Task
	row, err := tr.pool.Query(ctx, "SELECT id, user_id, title, description, is_completed, created_at, updated_at FROM tasks WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var t Task
		err := row.Scan(&t.Id,
			&t.UserId,
			&t.Title,
			&t.Description,
			&t.IsCompleted,
			&t.CreatedAt,
			&t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		return nil, ErrNoTasks
	}

	return tasks, nil
}

func (tr *TaskPgRepository) Create(ctx context.Context, task Task) (int, error) {
	var id int
	err := tr.pool.QueryRow(ctx, "INSERT INTO tasks (user_id, title, description) VALUES ($1, $2, $3) RETURNING id", task.UserId, task.Title, task.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (tr *TaskPgRepository) Delete(ctx context.Context, id int) error {
	_, err := tr.pool.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TaskPgRepository) UpdateTitle(ctx context.Context, newTitle string, id int) error {
	cmdTag, err := tr.pool.Exec(ctx, "UPDATE tasks SET title = $1, updated_at = $2 WHERE id = $3", newTitle, time.Now(), id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrTaskTitleNotUpdated
	}

	return nil
}

func (tr *TaskPgRepository) UpdateDescription(ctx context.Context, newDescription string, id int) error {
	cmdTag, err := tr.pool.Exec(ctx, "UPDATE tasks SET description = $1, updated_at = $2 WHERE id = $3", newDescription, time.Now(), id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrTaskDescNotUpdated
	}

	return nil
}

func (tr *TaskPgRepository) SwitchTaskStatus(ctx context.Context, id int) error {
	cmdTag, err := tr.pool.Exec(ctx, "UPDATE tasks SET is_completed = NOT is_completed, updated_at = $1 WHERE id = $2", time.Now(), id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrTaskStatusNotSwitched
	}

	return nil
}
