// +build integr_full

package main_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/common/catalog"
	"github.com/Vesninovich/go-tasks/book-store/common/orders"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	catalogservice "github.com/Vesninovich/go-tasks/book-store/orders/catalog/service"
	orderservice "github.com/Vesninovich/go-tasks/book-store/orders/order/service"
	ordersql "github.com/Vesninovich/go-tasks/book-store/orders/order/sql"
	"github.com/Vesninovich/go-tasks/book-store/orders/server"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const dbURL = "postgresql://gobookstoreorders@localhost:5432/gobookstore"
const catalogURL = "localhost:8001"
const schema = "orders_test"
const bufsize = 1024 * 1024

var lis *bufconn.Listener
var ctx = context.Background()
var client orders.OrdersClient

// TestMain tests
func TestMain(m *testing.M) {
	db, r := initSQL()

	cConn, err := grpc.Dial(catalogURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to catalog service on %s: %s", catalogURL, err)
	}

	cc := catalog.NewCatalogClient(cConn)
	c := catalogservice.New(cc)
	s := orderservice.New(r, c)

	lis = bufconn.Listen(bufsize)
	grpcServer := grpc.NewServer()
	orders.RegisterOrdersServer(grpcServer, server.New(s))
	log.Println("Starting gRPC server")
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %s", err)
		}
	}()

	oConn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %s", err)
	}
	client = orders.NewOrdersClient(oConn)

	res := m.Run()

	oConn.Close()
	cConn.Close()
	grpcServer.GracefulStop()
	db.Close()
	os.Exit(res)
}

func TestOrders(t *testing.T) {
	// aID, _ := uuid.FromString("15f76b4d-3ccb-417a-817c-c3a98e40ba34")
	bID, err := uuid.FromString("5d9644cd-6be5-41f2-97d3-15ecd509568a")
	if err != nil {
		t.Fatal(err)
	}
	created, err := client.CreateOrder(context.Background(), &orders.OrderCreateDTO{
		Description: "Test order",
		Book:        bID[:],
	})
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.GetOrder(context.Background(), &orders.OrderID{Id: created.Id})
	if err != nil {
		t.Fatal(err)
	}
	createdID, err := uuid.FromBytes(created.Id)
	if err != nil {
		t.Fatal(err)
	}
	resID, err := uuid.FromBytes(res.Id)
	if err != nil {
		t.Fatal(err)
	}
	createdBID, err := uuid.FromBytes(created.Book)
	if err != nil {
		t.Fatal(err)
	}
	resBID, err := uuid.FromBytes(res.Book)
	if err != nil {
		t.Fatal(err)
	}
	if resID != createdID || res.Description != created.Description || resBID != createdBID {
		t.Error("Created and queried orders are not equal")
	}
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func initSQL() (*sqlx.DB, *ordersql.Repository) {
	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB at URL %s\n%s", dbURL, err)
	}
	log.Printf("Connected to DB at URL %s\n", dbURL)

	log.Println("Creating schema")
	s := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schema)
	log.Println(s)
	db.MustExec(s)

	r := ordersql.New(db, schema)

	log.Println("Creating tables")
	log.Println(r.CreateTableStmt())
	db.MustExec(r.CreateTableStmt())
	log.Println("Finished setting up DB")

	return db, r
}
