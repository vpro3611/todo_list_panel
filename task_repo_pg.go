package main

import (
	"context"
	"errors"
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

	if tasks == nil || len(tasks) == 0 {
		return []Task{}, nil
	}

	return tasks, nil
}

func (tr *TaskPgRepository) GetByUserId(ctx context.Context, id int, actorID int, actorRole string) ([]Task, error) {
	var tasks []Task
	query := "SELECT id, user_id, title, description, is_completed, created_at, updated_at FROM tasks WHERE user_id = $1 AND (user_id = $2 OR $3 = 'admin')"
	row, err := tr.pool.Query(ctx, query, id, actorID, actorRole)
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

	if tasks == nil || len(tasks) == 0 {
		return []Task{}, nil
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

func (tr *TaskPgRepository) Delete(ctx context.Context, id int, actorId int, actorRole string) error {
	query := "DELETE FROM tasks WHERE id = $1 AND (user_id = $2 OR $3 = 'admin')"
	cmdTag, err := tr.pool.Exec(ctx, query, id, actorId, actorRole)

	if cmdTag.RowsAffected() == 0 {
		return errors.New("task not deleted")
	}

	if err != nil {
		return err
	}

	return nil
}

func (tr *TaskPgRepository) UpdateTitle(ctx context.Context, newTitle string, id int, actorId int, actorRole string) error {
	query := "UPDATE tasks SET title = $1, updated_at = $2 WHERE id = $3 AND (user_id = $4 OR $5 = 'admin')"
	cmdTag, err := tr.pool.Exec(ctx, query, newTitle, time.Now(), id, actorId, actorRole)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrTaskTitleNotUpdated
	}

	return nil
}

func (tr *TaskPgRepository) UpdateDescription(ctx context.Context, newDescription string, id int, actorId int, actorRole string) error {
	query := "UPDATE tasks SET description = $1, updated_at = $2 WHERE id = $3 AND (user_id = $4 OR $5 = 'admin')"
	cmdTag, err := tr.pool.Exec(ctx, query, newDescription, time.Now(), id, actorId, actorRole)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrTaskDescNotUpdated
	}

	return nil
}

func (tr *TaskPgRepository) SwitchTaskStatus(ctx context.Context, id int, actorId int, actorRole string) error {
	query := "UPDATE tasks SET is_completed = NOT is_completed, updated_at = $1 WHERE id = $2 AND (user_id = $3 OR $4 = 'admin')"
	cmdTag, err := tr.pool.Exec(ctx, query, time.Now(), id, actorId, actorRole)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrTaskStatusNotSwitched
	}

	return nil
}

func (tr *TaskPgRepository) GetTaskById(ctx context.Context, id int, actorId int, actorRole string) (*Task, error) {
	query := "SELECT id, user_id, title, description, is_completed, created_at, updated_at FROM tasks WHERE id = $1 AND (user_id = $2 OR $3 = 'admin')"

	var task Task

	err := tr.pool.QueryRow(ctx, query, id, actorId, actorRole).Scan(&task.Id,
		&task.UserId,
		&task.Title,
		&task.Description,
		&task.IsCompleted,
		&task.CreatedAt,
		&task.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &task, nil

	// type Task struct {
	//	Id          int       `json:"id"`
	//	UserId      int       `json:"user_id"`
	//	Title       string    `json:"title"`
	//	Description string    `json:"description"`
	//	IsCompleted bool      `json:"is_completed"`
	//	CreatedAt   time.Time `json:"created_at"`
	//	UpdatedAt   time.Time `json:"updated_at"`
	//}
}
