package taskservice

import (
	"context"
	"testing"
	"time"

	"github.com/Vesninovich/go-tasks/todos/common"
	"github.com/Vesninovich/go-tasks/todos/task"
	"github.com/Vesninovich/go-tasks/todos/task/inmemory"
)

func TestCreateValid(t *testing.T) {
	s := createService()
	created, err := s.CreateTask(context.Background(), "test", "", time.Now().Unix())
	if err != nil {
		t.Errorf("Got error while creating valid task: %s", err)
	}
	if created.Status != task.New {
		t.Errorf("Created task with incorrect status: %d", created.Status)
	}
}

func TestCreateWithEmptyName(t *testing.T) {
	s := createService()
	_, err := s.CreateTask(context.Background(), "", "", time.Now().Unix())
	if err == nil {
		t.Errorf("Expected to get error while creating task with empty name")
	}
	if _, ok := err.(*common.InvalidInputError); !ok {
		t.Errorf("Wrong error type from creating task with empty name")
	}
}

func TestCreateWithNegativeDueDate(t *testing.T) {
	s := createService()
	_, err := s.CreateTask(context.Background(), "test", "", -1234)
	if err == nil {
		t.Errorf("Expected to get error while creating task with negative due date")
	}
	if _, ok := err.(*common.InvalidInputError); !ok {
		t.Errorf("Wrong error type from creating task with negative due date")
	}
}

func createService() *Service {
	return New(inmemory.New())
}
