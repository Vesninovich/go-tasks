package tasksql

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Vesninovich/go-tasks/todos/common"
	"github.com/Vesninovich/go-tasks/todos/task"
)

type SQLRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLRepository {
	return &SQLRepository{db}
}

func (r *SQLRepository) Read(ctx context.Context, from, count uint) ([]task.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM tasks;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var t task.Task
	tasks := make([]task.Task, 0)
	for rows.Next() {
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.DueDate, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *SQLRepository) ReadOne(ctx context.Context, id uint64) (task.Task, error) {
	var t task.Task
	err := r.db.QueryRowContext(
		ctx, "SELECT * FROM tasks WHERE id=$1;", id,
	).Scan(&t.ID, &t.Name, &t.Description, &t.DueDate, &t.Status)
	if err == sql.ErrNoRows {
		return t, notFoundError(id)
	}
	return t, err
}

func (r *SQLRepository) Create(ctx context.Context, dto task.DTO) (task.Task, error) {
	var t task.Task
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO tasks (name, description, dueDate, status)
			VALUES ($1, $2, $3, $4)
			RETURNING id, name, description, dueDate, status;`,
		dto.Name, dto.Description, dto.DueDate, dto.Status,
	).Scan(&t.ID, &t.Name, &t.Description, &t.DueDate, &t.Status)
	return t, err
}

func (r *SQLRepository) Update(ctx context.Context, id uint64, dto task.DTO) (task.Task, error) {
	var t task.Task
	err := r.db.QueryRowContext(
		ctx,
		`UPDATE tasks
			SET name=$2, description=$3, dueDate=$4, status=$5
			WHERE id=$1
			RETURNING id, name, description, dueDate, status;`,
		id, dto.Name, dto.Description, dto.DueDate, dto.Status,
	).Scan(&t.ID, &t.Name, &t.Description, &t.DueDate, &t.Status)
	if err == sql.ErrNoRows {
		return t, notFoundError(id)
	}
	return t, err
}

func (r *SQLRepository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id=$1;", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if count == 0 {
		return notFoundError(id)
	}
	return err
}

func notFoundError(id uint64) *common.NotFoundError {
	return &common.NotFoundError{What: "Task with ID " + strconv.FormatUint(id, 10)}
}
