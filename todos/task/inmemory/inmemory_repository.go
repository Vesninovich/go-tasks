package inmemory

import (
	"context"
	"sync"

	"github.com/Vesninovich/go-tasks/todos/task"
)

type repository struct {
	tasks []task.Task
	lock  sync.RWMutex
	id    uint64
}

func New() *repository {
	return &repository{
		tasks: make([]task.Task, 0),
	}
}

func (r *repository) ReadAll(ctx context.Context) ([]task.Task, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.tasks, nil
}

func (r *repository) Create(ctx context.Context, taskDTO task.DTO) (task.Task, error) {
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
