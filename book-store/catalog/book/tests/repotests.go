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

var aut1, aut2 book.Author
var cat1, cat2 book.Category
var books = []bookrepo.CreateDTO{
	{Name: "bookA"},
	{Name: "bookB"},
}
var ctx = context.Background()

// Constructor is function that constructs repository to test
type Constructor func(t *testing.T) (author.Repository, category.Repository, bookrepo.Repository)

// RepoGet tests getting item with query
func RepoGet(t *testing.T, c Constructor) {
	repo := setup(t, c)
	count := uint(len(books))
	from := uint(1)
	res, err := repo.Get(ctx, 0, count, book.Query{})
	if err != nil {
		t.Fatalf("Error getting all books: %s", err)
	}
	if len(res) != len(books) {
		t.Fatalf("Expected to get all %d books, got %d", count, len(res))
	}
	stored := res

	t.Run("books from 1", func(t *testing.T) {
		res, err = repo.Get(ctx, from, count, book.Query{})
		if err != nil {
			t.Fatalf("Error getting books: %s", err)
		}
		if uint(len(res)) != count-from {
			t.Fatalf("Expected to get all books minus 1, got %d", len(res))
		}
	})

	t.Run("1 book", func(t *testing.T) {
		res, err = repo.Get(ctx, 0, 1, book.Query{})
		if err != nil {
			t.Fatalf("Error getting books: %s", err)
		}
		if len(res) != 1 {
			t.Fatalf("Expected to get 1 book, got %d", len(res))
		}
	})

	t.Run("book by author", func(t *testing.T) {
		res, err = repo.Get(ctx, 0, count, book.Query{Author: aut2.ID})
		if err != nil {
			t.Fatalf("Error getting books: %s", err)
		}
		if len(res) != 1 {
			t.Fatalf("Expected to get 1 book, got %d", len(res))
		}
		if res[0].ID != stored[1].ID {
			t.Error("Got wrong book with query by author")
		}
	})

	t.Run("book by categories", func(t *testing.T) {
		res, err = repo.Get(ctx, 0, count, book.Query{
			Categories: []uuid.UUID{cat2.ID},
		})
		if err != nil {
			t.Fatalf("Error getting books: %s", err)
		}
		if len(res) != 1 {
			t.Fatalf("Expected to get 1 book, got %d", len(res))
		}
		if res[0].ID != stored[1].ID {
			t.Error("Got wrong book with query by category")
		}
	})

	t.Run("book by author and categories", func(t *testing.T) {
		res, err = repo.Get(ctx, 0, count, book.Query{
			Author:     aut2.ID,
			Categories: []uuid.UUID{cat1.ID, cat2.ID},
		})
		if err != nil {
			t.Fatalf("Error getting books: %s", err)
		}
		if len(res) != 1 {
			t.Fatalf("Expected to get 1 book, got %d", len(res))
		}
		if res[0].ID != stored[1].ID {
			t.Error("Got wrong book with query by category")
		}
	})

	t.Run("book by id", func(t *testing.T) {
		res, err = repo.Get(ctx, 0, count, book.Query{
			ID: stored[1].ID,
		})
		if err != nil {
			t.Fatalf("Error getting books: %s", err)
		}
		if len(res) != 1 {
			t.Fatalf("Expected to get 1 book, got %d", len(res))
		}
		if res[0].ID != stored[1].ID ||
			res[0].Author.ID != books[1].Author.ID ||
			len(res[0].Categories) != len(books[1].Categories) {
			t.Error("Got wrong book with query by ID")
		}
	})

	t.Run("no books", func(t *testing.T) {
		res, err = repo.Get(ctx, 0, count, book.Query{
			Author: uuid.New(),
		})
		if err != nil {
			t.Fatalf("Error getting books: %s", err)
		}
		if len(res) != 0 {
			t.Fatalf("Expected to get 1 book, got %d", len(res))
		}
	})
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
	_, err := repo.Update(ctx, book.Book{
		ID:         id,
		Name:       name,
		Author:     aut2,
		Categories: []book.Category{cat1, cat2},
	})
	if err != nil {
		t.Fatalf("Error while updating item: %s", err)
	}
	b, err := repo.Get(ctx, 0, 1, book.Query{ID: id})
	if err != nil {
		t.Fatalf("Error while getting item: %s", err)
	}
	if len(b) != 1 {
		t.Fatalf("Failed to get updated item")
	}
	replaced := b[0]
	if replaced.Name != name ||
		replaced.Author.ID != aut2.ID ||
		len(replaced.Categories) != 2 {
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
	replaced, err := repo.Update(ctx, book.Book{
		ID:         item.ID,
		Name:       name,
		Author:     aut2,
		Categories: []book.Category{cat2},
	})
	if err != nil {
		t.Errorf("Error while updating item with some deleted: %s", err)
	}
	if replaced.Name != name ||
		replaced.Author.ID != aut2.ID ||
		replaced.Categories[0].ID != cat2.ID {
		t.Errorf("Expected to update item with data")
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
		b, err := repo.Get(ctx, 0, 1, book.Query{ID: item.ID})
		if err != nil {
			t.Fatalf("Error while getting item: %s", err)
		}
		if len(b) != 0 {
			t.Error("Did not expect to get deleted item")
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
	aut1, err = authorRepo.Create(ctx, author.CreateDTO{Name: "authorA"})
	if err != nil {
		t.Fatal("Error in setup: failed to create author")
	}
	aut2, err = authorRepo.Create(ctx, author.CreateDTO{Name: "authorB"})
	if err != nil {
		t.Fatal("Error in setup: failed to create author")
	}
	cat1, err = categoryRepo.Create(ctx, category.CreateDTO{Name: "catA"})
	if err != nil {
		t.Fatal("Error in setup: failed to create category")
	}
	cat2, err = categoryRepo.Create(ctx, category.CreateDTO{Name: "catB"})
	if err != nil {
		t.Fatal("Error in setup: failed to create category")
	}
	books[0].Author = aut1
	books[0].Categories = []book.Category{cat1}
	books[1].Author = aut2
	books[1].Categories = []book.Category{cat1, cat2}

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
