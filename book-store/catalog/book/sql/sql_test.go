// +build sql

package sql_test

import (
	"fmt"
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
const schema = "catalog_books_test"

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB at URL %s\n%s", dbURL, err)
	}
	defer db.Close()
	db.MustExec("CREATE SCHEMA IF NOT EXISTS " + schema)
	a := authorsql.New(db, schema)
	log.Println(a.CreateTableStmt())
	db.MustExec(a.CreateTableStmt())
	c := categorysql.New(db, schema)
	log.Println(c.CreateTableStmt())
	db.MustExec(c.CreateTableStmt())
	b := booksql.New(db, schema)
	log.Println(b.CreateTableStmt())
	db.MustExec(b.CreateTableStmt())
	res := m.Run()
	db.MustExec(fmt.Sprintf("DROP SCHEMA %s CASCADE;", schema))
	os.Exit(res)
}

func constructor(t *testing.T) (author.Repository, category.Repository, bookrepo.Repository) {
	t.Cleanup(clear)
	return authorsql.New(db, schema), categorysql.New(db, schema), booksql.New(db, schema)
}

func clear() {
	db.MustExec(fmt.Sprintf("DELETE FROM %s.books_categories;", schema))
	db.MustExec(fmt.Sprintf("DELETE FROM %s.books;", schema))
	db.MustExec(fmt.Sprintf("DELETE FROM %s.categories;", schema))
	db.MustExec(fmt.Sprintf("DELETE FROM %s.authors;", schema))
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
