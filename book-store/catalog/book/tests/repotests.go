package tests

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

var aut = book.Author{ID: uuid.New(), Name: "authorA"}
var cat = book.Category{ID: uuid.New(), Name: "catA"}
var books = []bookrepo.CreateDTO{
	{
		Name:       "bookA",
		Author:     aut,
		Categories: []book.Category{cat},
	},
	{
		Name: "bookB",
		Categories: []book.Category{
			cat,
			{},
		},
	},
}
var exAuthor book.Author
var exCategory book.Category
var ctx = context.Background()

// Constructor is function that constructs repository to test
type Constructor func(t *testing.T) (author.Repository, category.Repository, bookrepo.Repository)

// RepoGet tests getting item with query
func RepoGet(t *testing.T, c Constructor) {
	repo := setup(t, c)
	count := uint(len(books))
	res, err := repo.Get(ctx, 0, count, book.Query{})
	if err != nil {
		t.Errorf("Error getting all books: %s", err)
	}
	if len(res) != len(books) {
		t.Errorf("Expected to get all books, got only %d", len(res))
	}
	stored := res

	from := uint(1)
	res, err = repo.Get(ctx, from, count, book.Query{})
	if err != nil {
		t.Errorf("Error getting books: %s", err)
	}
	if uint(len(res)) != count-from {
		t.Errorf("Expected to get all books minus 1, got %d", len(res))
	}

	res, err = repo.Get(ctx, 0, 1, book.Query{})
	if err != nil {
		t.Errorf("Error getting books: %s", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected to get 1 book, got %d", len(res))
	}

	res, err = repo.Get(ctx, 0, count, book.Query{Author: exAuthor.ID})
	if err != nil {
		t.Errorf("Error getting books: %s", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected to get 1 book, got %d", len(res))
	}
	if res[0].ID != stored[1].ID {
		t.Error("Got wrong book with query by author")
	}

	res, err = repo.Get(ctx, 0, count, book.Query{
		Categories: []uuid.UUID{exCategory.ID},
	})
	if err != nil {
		t.Errorf("Error getting books: %s", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected to get 1 book, got %d", len(res))
	}
	if res[0].ID != stored[1].ID {
		t.Error("Got wrong book with query by category")
	}

	res, err = repo.Get(ctx, 0, count, book.Query{
		Author:     aut.ID,
		Categories: []uuid.UUID{exCategory.ID},
	})
	if err != nil {
		t.Errorf("Error getting books: %s", err)
	}
	if len(res) != 0 {
		t.Errorf("Expected to get no books, got %d", len(res))
	}

	res, err = repo.Get(ctx, 0, count, book.Query{
		ID: stored[1].ID,
	})
	if err != nil {
		t.Errorf("Error getting books: %s", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected to get 1 book, got %d", len(res))
	}
	if res[0].ID != stored[1].ID {
		t.Error("Got wrong book with query by ID")
	}
}

// RepoCreate tests creating items
func RepoCreate(t *testing.T, c Constructor) {
	setup(t, c)
}

// RepoUpdate tests updating item
func RepoUpdate(t *testing.T, c Constructor) {
	repo, stored := setupMutation(t, c)
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

// RepoUpdateNonExisting tests updating non-existing item
func RepoUpdateNonExisting(t *testing.T, c Constructor) {
	repo := setup(t, c)
	_, err := repo.Update(ctx, book.Book{
		ID:         uuid.New(),
		Name:       "",
		Author:     book.Author{},
		Categories: []book.Category{},
	})
	checkNotFound(t, err)
}

// RepoUpdateDeleted tests updating deleted item
func RepoUpdateDeleted(t *testing.T, c Constructor) {
	repo, id, _ := setupAlreadyDeleted(t, c)
	_, err := repo.Update(ctx, book.Book{
		ID:         id,
		Name:       "",
		Author:     book.Author{},
		Categories: []book.Category{},
	})
	checkNotFound(t, err)
}

// RepoUpdateWithSomeDeleted tests updating item if some are deleted
func RepoUpdateWithSomeDeleted(t *testing.T, c Constructor) {
	repo, id, stored := setupAlreadyDeleted(t, c)
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

// RepoDelete tests deleting item
func RepoDelete(t *testing.T, c Constructor) {
	repo, stored := setupMutation(t, c)
	for _, item := range stored {
		_, err := repo.Delete(ctx, item.ID)
		if err != nil {
			t.Errorf("Error while deleting item: %s", err)
		}
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

func findByName(name string, data []book.Book, t *testing.T) book.Book {
	for _, item := range data {
		if item.Name == name {
			return item
		}
	}
	t.Errorf("Item with name %s not found", name)
	return book.Book{}
}

func setup(t *testing.T, c Constructor) bookrepo.Repository {
	authorRepo, categoryRepo, repo := c(t)

	var err error
	exAuthor, err = authorRepo.Create(ctx, author.CreateDTO{Name: "authorB"})
	if err != nil {
		t.Fatal("Error in setup: failed to create author")
	}
	exCategory, err = categoryRepo.Create(ctx, category.CreateDTO{Name: "catB"})
	if err != nil {
		t.Fatal("Error in setup: failed to create category")
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
			t.Fatalf("Error while creating category %s: %s", a, err)
		}
	}
	return repo
}

func setupMutation(t *testing.T, c Constructor) (bookrepo.Repository, []book.Book) {
	repo := setup(t, c)
	stored, err := repo.Get(ctx, 0, uint(len(books)), book.Query{})
	if err != nil {
		t.Fatalf("Error while fetching all data: %s", err)
	}
	return repo, stored
}

func setupAlreadyDeleted(t *testing.T, c Constructor) (bookrepo.Repository, uuid.UUID, []book.Book) {
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
		t.Errorf("Expected to get NotFound error")
	}
	if _, typeCorrect := err.(*commonerrors.NotFound); !typeCorrect {
		t.Errorf("Expected to get error of type *commonerrors.NotFound, got %T", err)
	}
}
