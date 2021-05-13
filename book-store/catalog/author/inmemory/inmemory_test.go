package inmemory_test

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	"github.com/Vesninovich/go-tasks/book-store/catalog/author/inmemory"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

var authors = []author.CreateDTO{
	{Name: "Л.Н. Толстой"},
	{Name: "Ф.М. Достоевский"},
}
var ctx = context.Background()

func TestGetAll(t *testing.T) {
	repo := setup(t)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Error while getting all stored items: %s", err)
	}
	if len(stored) != len(authors) {
		t.Errorf("Expected to have %d items stored, got %d", len(authors), len(stored))
	}
}

func TestGet(t *testing.T) {
	repo := setup(t)
	stored, _ := repo.GetAll(ctx)
	for _, item := range stored {
		found, err := repo.Get(ctx, item.ID)
		if err != nil {
			t.Errorf("Error while getting item: %s", err)
		}
		if found.Name != item.Name || found.ID != item.ID {
			t.Error("Got wrong item")
		}
	}
}

func TestGetNonExisting(t *testing.T) {
	repo := setup(t)
	_, err := repo.Get(ctx, uuid.New())
	checkNotFound(t, err)
}

func TestCreate(t *testing.T) {
	repo := setup(t)
	repo.GetAll(ctx)
}

func TestUpdate(t *testing.T) {
	repo, stored := setupMutation(t)
	id := stored[0].ID
	name := "И.С. Тургенев"
	replaced, err := repo.Update(ctx, author.UpdateDTO{
		ID:   id,
		Name: name,
	})
	if err != nil {
		t.Errorf("Error while updating item: %s", err)
	}
	if replaced.Name != name {
		t.Errorf("Expected to update item with data %s, got %s", name, replaced.Name)
	}
}

func TestUpdateNonExisting(t *testing.T) {
	repo := setup(t)
	_, err := repo.Update(ctx, author.UpdateDTO{
		ID:   uuid.New(),
		Name: "",
	})
	checkNotFound(t, err)
}

func TestUpdateDeleted(t *testing.T) {
	repo, id, _ := setupAlreadyDeleted(t)
	_, err := repo.Update(ctx, author.UpdateDTO{
		ID:   id,
		Name: "",
	})
	checkNotFound(t, err)
}

func TestUpdateWithSomeDeleted(t *testing.T) {
	repo, id, stored := setupAlreadyDeleted(t)
	var item book.Author
	for _, item = range stored {
		if item.ID != id {
			break
		}
	}
	name := "Hemingway"
	replaced, err := repo.Update(ctx, author.UpdateDTO{
		ID:   item.ID,
		Name: name,
	})
	if err != nil {
		t.Errorf("Error while updating item with some deleted: %s", err)
	}
	if replaced.Name != name {
		t.Errorf("Expected to update item with data %s, got %s", name, replaced.Name)
	}
}

func TestDelete(t *testing.T) {
	repo, stored := setupMutation(t)
	for _, item := range stored {
		_, err := repo.Delete(ctx, item.ID)
		if err != nil {
			t.Errorf("Error while deleting item: %s", err)
		}
	}
}

func TestDeleteTwice(t *testing.T) {
	repo, id, _ := setupAlreadyDeleted(t)
	_, err := repo.Delete(ctx, id)
	checkNotFound(t, err)
}

func TestDeleteNonExisting(t *testing.T) {
	repo := setup(t)
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

func setup(t *testing.T) author.Repository {
	repo := inmemory.New()
	for _, a := range authors {
		_, err := repo.Create(ctx, a)
		if err != nil {
			t.Errorf("Error while creating author %s: %s", a, err)
		}
	}
	return repo
}

func setupMutation(t *testing.T) (author.Repository, []book.Author) {
	repo := setup(t)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Error while fetching all data: %s", err)
	}
	return repo, stored
}

func setupAlreadyDeleted(t *testing.T) (author.Repository, uuid.UUID, []book.Author) {
	repo, stored := setupMutation(t)
	id := stored[0].ID
	_, err := repo.Delete(ctx, id)
	if err != nil {
		t.Errorf("Error deleting item: %s", err)
	}
	return repo, id, stored
}

func checkNotFound(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected to get NotFound error")
	}
	if _, typeCorrect := err.(*commonerrors.NotFound); !typeCorrect {
		t.Errorf("Expected to get error of type *commonerrors.NotFound, got %T", err)
	}
}
