package inmemory_test

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	authorInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/author/inmemory"
	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	"github.com/Vesninovich/go-tasks/book-store/catalog/book/inmemory"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	categoryInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/category/inmemory"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

var books = []bookrepo.CreateDTO{
	{
		Name:   "bookA",
		Author: book.Author{Name: "authorA"},
		Categories: []book.Category{
			{Name: "catA"},
		},
	},
	{
		Name: "bookB",
		Categories: []book.Category{
			{Name: "catA"},
			{},
		},
	},
}
var exAuthor book.Author
var exCategory book.Category
var ctx = context.Background()

func TestGetAll(t *testing.T) {
	repo := setup(t)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Error while getting all stored items: %s", err)
	}
	if len(stored) != len(books) {
		t.Errorf("Expected to have %d items stored, got %d", len(books), len(stored))
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
	name := "asddsa"
	aut := book.Author{Name: "authorC"}
	cat := []book.Category{{Name: "catC"}}
	replaced, err := repo.Update(ctx, book.Book{
		ID:         id,
		Name:       name,
		Author:     aut,
		Categories: cat,
	})
	if err != nil {
		t.Errorf("Error while updating item: %s", err)
	}
	if replaced.Name != name ||
		replaced.Author.Name != aut.Name ||
		replaced.Categories[0].Name != cat[0].Name {
		t.Errorf("Expected to update item with data")
	}
}

func TestUpdateNonExisting(t *testing.T) {
	repo := setup(t)
	_, err := repo.Update(ctx, book.Book{
		ID:         uuid.New(),
		Name:       "",
		Author:     book.Author{},
		Categories: []book.Category{},
	})
	checkNotFound(t, err)
}

func TestUpdateDeleted(t *testing.T) {
	repo, id, _ := setupAlreadyDeleted(t)
	_, err := repo.Update(ctx, book.Book{
		ID:         id,
		Name:       "",
		Author:     book.Author{},
		Categories: []book.Category{},
	})
	checkNotFound(t, err)
}

func TestUpdateWithSomeDeleted(t *testing.T) {
	repo, id, stored := setupAlreadyDeleted(t)
	var item book.Book
	for _, item = range stored {
		if item.ID != id {
			break
		}
	}
	name := "asddsa"
	aut := book.Author{Name: "authorC"}
	cat := []book.Category{{Name: "catC"}}
	replaced, err := repo.Update(ctx, book.Book{
		ID:         item.ID,
		Name:       name,
		Author:     aut,
		Categories: cat,
	})
	if err != nil {
		t.Errorf("Error while updating item with some deleted: %s", err)
	}
	if replaced.Name != name ||
		replaced.Author.Name != aut.Name ||
		replaced.Categories[0].Name != cat[0].Name {
		t.Errorf("Expected to update item with data")
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

func findByName(name string, data []book.Book, t *testing.T) book.Book {
	for _, item := range data {
		if item.Name == name {
			return item
		}
	}
	t.Errorf("Item with name %s not found", name)
	return book.Book{}
}

func setup(t *testing.T) bookrepo.Repository {
	authorRepo := authorInMemory.New()
	categoryRepo := categoryInMemory.New()
	repo := inmemory.New(authorRepo, categoryRepo)

	var err error
	exAuthor, err = authorRepo.Create(ctx, author.CreateDTO{Name: "authorB"})
	if err != nil {
		t.Error("Error in setup: failed to create author")
	}
	exCategory, err = categoryRepo.Create(ctx, category.CreateDTO{Name: "catB"})
	if err != nil {
		t.Error("Error in setup: failed to create category")
	}
	books[1].Author = book.Author{
		ID:   exAuthor.ID,
		Name: exAuthor.Name,
	}
	books[1].Categories[1] = book.Category{
		ID:   exCategory.ID,
		Name: exCategory.Name,
	}

	for _, a := range books {
		_, err := repo.Create(ctx, a)
		if err != nil {
			t.Errorf("Error while creating category %s: %s", a, err)
		}
	}
	return repo
}

func setupMutation(t *testing.T) (bookrepo.Repository, []book.Book) {
	repo := setup(t)
	stored, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Error while fetching all data: %s", err)
	}
	return repo, stored
}

func setupAlreadyDeleted(t *testing.T) (bookrepo.Repository, uuid.UUID, []book.Book) {
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
