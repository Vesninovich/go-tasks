// +build sql

package sql_test

import (
	"log"
	"os"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	authorsql "github.com/Vesninovich/go-tasks/book-store/catalog/author/sql"
	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	booksql "github.com/Vesninovich/go-tasks/book-store/catalog/book/sql"
	"github.com/Vesninovich/go-tasks/book-store/catalog/book/tests"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	categorysql "github.com/Vesninovich/go-tasks/book-store/catalog/category/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const dbURL = "postgresql://gobookstorecatalog@localhost:5432/gobookstore"

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB at URL %s\n%s", dbURL, err)
	}
	defer db.Close()
	log.Println(authorsql.Table)
	db.MustExec(authorsql.Table)
	log.Println(categorysql.Table)
	db.MustExec(categorysql.Table)
	log.Println(booksql.Table)
	db.MustExec(booksql.Table)
	res := m.Run()
	db.MustExec("DROP TABLE books_categories;")
	db.MustExec("DROP TABLE books;")
	db.MustExec("DROP TABLE categories;")
	db.MustExec("DROP TABLE authors;")
	os.Exit(res)
}

func constructor(t *testing.T) (author.Repository, category.Repository, bookrepo.Repository) {
	t.Cleanup(clear)
	return authorsql.New(db), categorysql.New(db), booksql.New(db)
}

func clear() {
	db.MustExec("DELETE FROM books_categories;")
	db.MustExec("DELETE FROM books;")
	db.MustExec("DELETE FROM categories;")
	db.MustExec("DELETE FROM authors;")
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
