package inmemory_test

import (
	"context"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	authorInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/author/inmemory"
	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	"github.com/Vesninovich/go-tasks/book-store/catalog/book/inmemory"
	"github.com/Vesninovich/go-tasks/book-store/catalog/book/tests"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	categoryInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/category/inmemory"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
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

func constructor(t *testing.T) (author.Repository, category.Repository, bookrepo.Repository) {
	return authorInMemory.New(), categoryInMemory.New(), inmemory.New()
}

func TestGet(t *testing.T) {
	tests.RepoGet(t, constructor)
}

func TestCreate(t *testing.T) {
	tests.RepoCreate(t, constructor)
}

func TestUpdate(t *testing.T) {
	tests.RepoUpdate(t, constructor)
}

func TestUpdateNonExisting(t *testing.T) {
	tests.RepoUpdateNonExisting(t, constructor)
}

func TestUpdateDeleted(t *testing.T) {
	tests.RepoUpdateDeleted(t, constructor)
}

func TestUpdateWithSomeDeleted(t *testing.T) {
	tests.RepoUpdateWithSomeDeleted(t, constructor)
}

func TestDelete(t *testing.T) {
	tests.RepoDelete(t, constructor)
}

func TestDeleteTwice(t *testing.T) {
	tests.RepoDeleteTwice(t, constructor)
}

func TestDeleteNonExisting(t *testing.T) {
	tests.RepoDeleteNonExisting(t, constructor)
}
