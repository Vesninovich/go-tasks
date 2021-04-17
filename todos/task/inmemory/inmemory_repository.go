package inmemory

import (
	"context"
	"sync"

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

// ReadAll reads all saved tasks
func (r *Repository) ReadAll(ctx context.Context) ([]task.Task, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.tasks, nil
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
