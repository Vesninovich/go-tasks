package inmemory_test

import (
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category/inmemory"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category/tests"
)

func constructor(t *testing.T) category.Repository {
	return inmemory.New()
}

func TestGetAll(t *testing.T) {
	tests.RepoGetAll(t, constructor)
}

func TestGet(t *testing.T) {
	tests.RepoGet(t, constructor)
}

func TestGetNonExisting(t *testing.T) {
	tests.RepoGetNonExisting(t, constructor)
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
