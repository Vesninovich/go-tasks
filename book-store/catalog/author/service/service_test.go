package authorservice_test

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author/inmemory"
	authorservice "github.com/Vesninovich/go-tasks/book-store/catalog/author/service"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
)

func TestCreateValid(t *testing.T) {
	s := createService()
	_, err := s.CreateAuthor(context.Background(), "test")
	if err != nil {
		t.Errorf("Got error while creating valid author: %s", err)
	}
}

func TestCreateWithEmptyName(t *testing.T) {
	s := createService()
	_, err := s.CreateAuthor(context.Background(), "")
	if err == nil {
		t.Error("Expected to get error while creating author with empty name")
	}
	if _, ok := err.(*commonerrors.InvalidInput); !ok {
		t.Errorf("Wrong error type from creating author with empty name, got %T", err)
	}
}

func createService() *authorservice.Service {
	return authorservice.New(inmemory.New())
}
