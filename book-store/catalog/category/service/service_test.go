package categoryservice_test

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/category/inmemory"
	categoryservice "github.com/Vesninovich/go-tasks/book-store/catalog/category/service"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
)

func TestCreateValid(t *testing.T) {
	s := createService()
	_, err := s.CreateCategory(context.Background(), "test")
	if err != nil {
		t.Errorf("Got error while creating valid category: %s", err)
	}
}

func TestCreateWithEmptyName(t *testing.T) {
	s := createService()
	_, err := s.CreateCategory(context.Background(), "")
	if err == nil {
		t.Error("Expected to get error while creating category with empty name")
	}
	if _, ok := err.(*commonerrors.InvalidInput); !ok {
		t.Errorf("Wrong error type from creating category with empty name, got %T", err)
	}
}

func createService() *categoryservice.Service {
	return categoryservice.New(inmemory.New())
}
