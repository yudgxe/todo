package dao

import (
	"context"
	"errors"
	"fmt"
	"medods/database"
	"medods/database/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// todo: warp in error packge
type ErrTaskNotFound struct {
	Id int32
}

func NewErrTaskNotFound(id int32) *ErrTaskNotFound {
	return &ErrTaskNotFound{Id: id}
}

func (this ErrTaskNotFound) Error() string {
	return fmt.Sprintf("task with id: %v not found", this.Id)
}

var tasks TasksDao

type TasksDao interface {
	Create(ctx context.Context, task *model.Task) error
	List(ctx context.Context) ([]model.Task, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, task *model.Task) error
}

type tasksDao struct {
	db *pgxpool.Pool
}

func Tasks() TasksDao {
	if tasks == nil {
		tasks = &tasksDao{db: database.GetDatabase()}
	}

	return tasks
}

func (this tasksDao) Create(ctx context.Context, task *model.Task) error {
	return this.db.QueryRow(ctx,
		"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id, status, created_at",
		task.Title, task.Description, task.Status,
	).Scan(&task.Id, &task.Status, &task.CreatedAt)
}

func (this tasksDao) List(ctx context.Context) ([]model.Task, error) {
	rows, err := this.db.Query(ctx, "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	var task model.Task

	for rows.Next() {
		if err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (this tasksDao) Delete(ctx context.Context, id int32) error {
	result, err := this.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)

	if result.RowsAffected() == 0 {
		return NewErrTaskNotFound(id)
	}

	return err
}

func (this tasksDao) Update(ctx context.Context, task *model.Task) error {
	err := this.db.QueryRow(ctx,
		"UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5 RETURNING created_at, updated_at",
		task.Title, task.Description, task.Status, task.UpdatedAt, task.Id,
	).Scan(&task.CreatedAt, &task.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return NewErrTaskNotFound(task.Id)
	}

	return err
}
