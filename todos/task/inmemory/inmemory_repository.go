package inmemory

import (
	"context"
	"strconv"
	"sync"

	"github.com/Vesninovich/go-tasks/todos/common"
	"github.com/Vesninovich/go-tasks/todos/task"
)

// Repository represents an in-memory repository of tasks
type Repository struct {
	tasks []task.Task
	lock  sync.RWMutex
	id    uint64
}

// New creates new instance of in-memory Repository
func New() *Repository {
	return &Repository{
		tasks: make([]task.Task, 0),
	}
}

// Read reads `count` saved tasks starting from `from`
func (r *Repository) Read(ctx context.Context, from, count uint) ([]task.Task, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if from > uint(len(r.tasks)-1) {
		return make([]task.Task, 0), nil
	}

	if count == 0 {
		return r.tasks[from:], nil
	}

	to := int(from + count)
	if to > len(r.tasks) {
		to = len(r.tasks)
	}
	return r.tasks[from:to], nil
}

// ReadOne searches for task with given id, returns error if it is not found
func (r *Repository) ReadOne(ctx context.Context, id uint64) (task.Task, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var empty task.Task
	for _, tsk := range r.tasks {
		if id == tsk.ID {
			return tsk, nil
		}
	}
	return empty, &common.NotFoundError{What: "Task with ID " + strconv.FormatUint(id, 10)}
}

// Create adds new task to saved
func (r *Repository) Create(ctx context.Context, taskDTO task.DTO) (task.Task, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	t := task.Task{
		ID:          r.id,
		Name:        taskDTO.Name,
		Description: taskDTO.Description,
		DueDate:     taskDTO.DueDate,
		Status:      taskDTO.Status,
	}
	r.id++

	r.tasks = append(r.tasks, t)
	return t, nil
}

// Update updates task with given id, returns error if it is not found
func (r *Repository) Update(ctx context.Context, id uint64, taskDTO task.DTO) (task.Task, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, tsk := range r.tasks {
		if id == tsk.ID {
			update := task.Task{
				ID:          id,
				Name:        taskDTO.Name,
				Description: taskDTO.Description,
				DueDate:     taskDTO.DueDate,
				Status:      taskDTO.Status,
			}
			r.tasks[i] = update
			return update, nil
		}
	}
	var empty task.Task
	return empty, notFoundError(id)
}

// Delete deletes task with given id, returns error if it is not found
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, tsk := range r.tasks {
		if id == tsk.ID {
			r.tasks = append(r.tasks[0:i], r.tasks[i+1:len(r.tasks)]...)
			return nil
		}
	}
	return notFoundError(id)
}

func notFoundError(id uint64) *common.NotFoundError {
	return &common.NotFoundError{What: "Task with ID " + strconv.FormatUint(id, 10)}
}
