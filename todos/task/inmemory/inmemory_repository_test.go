package inmemory_test

import (
	"context"
	"testing"
	"time"

	"github.com/Vesninovich/go-tasks/todos/task"
	"github.com/Vesninovich/go-tasks/todos/task/inmemory"
)

var tasks = []task.DTO{
	{
		Name:        "testA",
		Description: "",
		DueDate:     time.Now(),
		Status:      task.New,
	},
	{
		Name:        "testB",
		Description: "some description",
		DueDate:     time.Now().Add(time.Duration(24 * 60 * 60 * 1000 * 1000)),
		Status:      task.InProgress,
	},
}

func TestCreate(t *testing.T) {
	r := inmemory.New()
	t0, err := r.Create(context.Background(), tasks[0])
	if err != nil {
		t.Errorf("Got error while saving task 0: %s", err)
	}
	t1, err := r.Create(context.Background(), tasks[1])
	if err != nil {
		t.Errorf("Got error while saving task 1: %s", err)
	}
	saved, err := r.ReadAll(context.Background())
	switch {
	case err != nil:
		t.Errorf("Got error while reading tasks: %s", err)
	case len(saved) != 2:
		t.Errorf("Got wrong number of saved tasks: %d", len(saved))
	case t0.ID == t1.ID:
		t.Errorf("Saved tasks have same ID")
	}
}
