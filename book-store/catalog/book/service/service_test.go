package bookservice_test

import (
	"context"
	"testing"

	authorInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/author/inmemory"
	authorservice "github.com/Vesninovich/go-tasks/book-store/catalog/author/service"
	"github.com/Vesninovich/go-tasks/book-store/catalog/book/inmemory"
	bookservice "github.com/Vesninovich/go-tasks/book-store/catalog/book/service"
	categoryInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/category/inmemory"
	categoryservice "github.com/Vesninovich/go-tasks/book-store/catalog/category/service"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

var ctx = context.Background()

var author = book.Author{
	Name: "Test author",
}
var categories = []book.Category{
	{Name: "Test1"},
	{Name: "Test2"},
}

func TestCreate(t *testing.T) {
	// TODO: test nested creation
	s := setup(t)
	name := "Test"
	res, err := s.CreateBook(ctx, name, author, categories)
	if err != nil {
		t.Errorf("Error while creating valid book: %s", err)
	}
	if res.Name != name {
		t.Errorf("Wrong name: got %s", res.Name)
	}
}

func TestCreateInvalidName(t *testing.T) {
	s := setup(t)
	name := ""
	_, err := s.CreateBook(ctx, name, author, categories)
	if err == nil {
		t.Error("Expected to get error for empty name")
	}
	if _, ok := err.(*commonerrors.InvalidInput); !ok {
		t.Errorf("Expected to get error of invalid input type, got %T", err)
	}
}

func TestCreateNonExistingFields(t *testing.T) {
	// TODO: update for nested creation
	t.Skip()
	s := setup(t)
	name := "Test"

	_, err := s.CreateBook(ctx, name, book.Author{ID: uuid.New()}, categories)
	if err == nil {
		t.Error("Expected to get error for non-existing author")
	}
	if _, ok := err.(*commonerrors.NotFound); !ok {
		t.Errorf("Expected to get error of not found type, got %T", err)
	}

	categories[1] = book.Category{ID: uuid.New()}
	_, err = s.CreateBook(ctx, name, author, categories)
	if err == nil {
		t.Error("Expected to get error for non-existing category")
	}
	if _, ok := err.(*commonerrors.NotFound); !ok {
		t.Errorf("Expected to get error of not found type, got %T", err)
	}
}

func setup(t *testing.T) *bookservice.BookService {
	as := authorservice.New(authorInMemory.New())
	cs := categoryservice.New(categoryInMemory.New())
	var err error
	author, err = as.CreateAuthor(ctx, author.Name)
	if err != nil {
		t.Fatalf("Error while creating author: %s", err)
	}
	for i, c := range categories {
		created, err := cs.CreateCategory(ctx, c.Name, c.ParentID)
		if err != nil {
			t.Fatalf("Error while creating category: %s", err)
		}
		categories[i] = created
	}
	return bookservice.New(inmemory.New(), as, cs)
}
