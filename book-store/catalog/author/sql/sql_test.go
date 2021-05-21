// +build sql

package sql_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	authorsql "github.com/Vesninovich/go-tasks/book-store/catalog/author/sql"
	"github.com/Vesninovich/go-tasks/book-store/catalog/author/tests"
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
	db.MustExec(authorsql.Table)
	res := m.Run()
	db.MustExec("DROP TABLE authors;")
	os.Exit(res)
}

func constructor(t *testing.T) author.Repository {
	t.Cleanup(clear)
	return authorsql.New(db)
}

func clear() {
	db.MustExecContext(context.Background(), "DELETE FROM authors;")
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
