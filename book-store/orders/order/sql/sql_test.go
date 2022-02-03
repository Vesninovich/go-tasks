// +build sql

package sql_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/orders/order"
	ordersql "github.com/Vesninovich/go-tasks/book-store/orders/order/sql"
	"github.com/Vesninovich/go-tasks/book-store/orders/order/tests"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const dbURL = "postgresql://gobookstoreorders@localhost:5432/gobookstore"
const schema = "orders_test"

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB at URL %s\n%s", dbURL, err)
	}
	defer db.Close()
	db.MustExec("CREATE SCHEMA IF NOT EXISTS " + schema)
	o := ordersql.New(db, schema)
	log.Println(o.CreateTableStmt())
	db.MustExec(o.CreateTableStmt())
	res := m.Run()
	db.MustExec(fmt.Sprintf("DROP SCHEMA %s CASCADE;", schema))
	os.Exit(res)
}

func constructor(t *testing.T) order.Repository {
	t.Cleanup(clear)
	return ordersql.New(db, schema)
}

func clear() {
	db.MustExecContext(context.Background(), fmt.Sprintf("DELETE FROM %s.orders;", schema))
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
