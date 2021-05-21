package tests

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

var authors = []author.CreateDTO{
	{Name: "Л.Н. Толстой"},
	{Name: "Ф.М. Достоевский"},
}
var ctx = context.Background()

// Constructor is function that constructs repository to test
type Constructor func(*testing.T) author.Repository

// RepoGetAll tests getting all data
func RepoGetAll(t *testing.T, c Constructor) {
	repo := setup(t, c)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Error while getting all stored items: %s", err)
	}
	if len(stored) != len(authors) {
		t.Errorf("Expected to have %d items stored, got %d", len(authors), len(stored))
	}
}

// RepoGet tests getting item by id
func RepoGet(t *testing.T, c Constructor) {
	repo := setup(t, c)
	stored, _ := repo.GetAll(ctx)
	for _, item := range stored {
		found, err := repo.Get(ctx, item.ID)
		if err != nil {
			t.Fatalf("Error while getting item: %s", err)
		}
		if found.Name != item.Name || found.ID != item.ID {
			t.Error("Got wrong item")
		}
	}
}

// RepoGetNonExisting tests getting non-existing item by id
func RepoGetNonExisting(t *testing.T, c Constructor) {
	repo := setup(t, c)
	_, err := repo.Get(ctx, uuid.New())
	checkNotFound(t, err)
}

// RepoCreate tests creating items
func RepoCreate(t *testing.T, c Constructor) {
	repo := setup(t, c)
	repo.GetAll(ctx)
}

// RepoUpdate tests updating item
func RepoUpdate(t *testing.T, c Constructor) {
	repo, stored := setupMutation(t, c)
	id := stored[0].ID
	name := "И.С. Тургенев"
	replaced, err := repo.Update(ctx, book.Author{
		ID:   id,
		Name: name,
	})
	if err != nil {
		t.Fatalf("Error while updating item: %s", err)
	}
	if replaced.Name != name {
		t.Errorf("Expected to update item with data %s, got %s", name, replaced.Name)
	}
}

// RepoUpdateNonExisting tests updating non-existing item
func RepoUpdateNonExisting(t *testing.T, c Constructor) {
	repo := setup(t, c)
	_, err := repo.Update(ctx, book.Author{
		ID:   uuid.New(),
		Name: "",
	})
	checkNotFound(t, err)
}

// RepoUpdateDeleted tests updating deleted item
func RepoUpdateDeleted(t *testing.T, c Constructor) {
	repo, id, _ := setupAlreadyDeleted(t, c)
	_, err := repo.Update(ctx, book.Author{
		ID:   id,
		Name: "",
	})
	checkNotFound(t, err)
}

// RepoUpdateWithSomeDeleted tests updating item if some are deleted
func RepoUpdateWithSomeDeleted(t *testing.T, c Constructor) {
	repo, id, stored := setupAlreadyDeleted(t, c)
	var item book.Author
	for _, item = range stored {
		if item.ID != id {
			break
		}
	}
	name := "Hemingway"
	replaced, err := repo.Update(ctx, book.Author{
		ID:   item.ID,
		Name: name,
	})
	if err != nil {
		t.Fatalf("Error while updating item with some deleted: %s", err)
	}
	if replaced.Name != name {
		t.Errorf("Expected to update item with data %s, got %s", name, replaced.Name)
	}
}

// RepoDelete tests deleting item
func RepoDelete(t *testing.T, c Constructor) {
	repo, stored := setupMutation(t, c)
	for _, item := range stored {
		_, err := repo.Delete(ctx, item.ID)
		if err != nil {
			t.Fatalf("Error while deleting item: %s", err)
		}
		_, err = repo.Get(ctx, item.ID)
		checkNotFound(t, err)
	}
}

// RepoDeleteTwice tests deleting item twice
func RepoDeleteTwice(t *testing.T, c Constructor) {
	repo, id, _ := setupAlreadyDeleted(t, c)
	_, err := repo.Delete(ctx, id)
	checkNotFound(t, err)
}

// RepoDeleteNonExisting tests deleting non-existing item
func RepoDeleteNonExisting(t *testing.T, c Constructor) {
	repo := setup(t, c)
	_, err := repo.Delete(ctx, uuid.New())
	checkNotFound(t, err)
}

func findByName(name string, data []book.Author, t *testing.T) book.Author {
	for _, item := range data {
		if item.Name == name {
			return item
		}
	}
	t.Errorf("Item with name %s not found", name)
	return book.Author{}
}

func setup(t *testing.T, c Constructor) author.Repository {
	repo := c(t)
	for _, a := range authors {
		_, err := repo.Create(ctx, a)
		if err != nil {
			t.Fatalf("Error while creating author %s: %s", a, err)
		}
	}
	return repo
}

func setupMutation(t *testing.T, c Constructor) (author.Repository, []book.Author) {
	repo := setup(t, c)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Error while fetching all data: %s", err)
	}
	return repo, stored
}

func setupAlreadyDeleted(t *testing.T, c Constructor) (author.Repository, uuid.UUID, []book.Author) {
	repo, stored := setupMutation(t, c)
	id := stored[0].ID
	_, err := repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Error deleting item: %s", err)
	}
	return repo, id, stored
}

func checkNotFound(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Expected to get NotFound error")
	}
	if _, typeCorrect := err.(*commonerrors.NotFound); !typeCorrect {
		t.Errorf("Expected to get error of type *commonerrors.NotFound, got %T", err)
	}
}
