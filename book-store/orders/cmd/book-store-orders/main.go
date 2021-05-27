package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Vesninovich/go-tasks/book-store/common/catalog"
	"github.com/Vesninovich/go-tasks/book-store/common/orders"
	catalogservice "github.com/Vesninovich/go-tasks/book-store/orders/catalog/service"
	orderservice "github.com/Vesninovich/go-tasks/book-store/orders/order/service"
	ordersql "github.com/Vesninovich/go-tasks/book-store/orders/order/sql"
	"github.com/Vesninovich/go-tasks/book-store/orders/server"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

// @title Book Store Orders Service
// @version 0.0
// @description Service for placing and reading book orders

// @contact.name Dimas
// @contact.url https://github.com/Vesninovich
// @contact.email dmitry@vesnin.work

// @license.name ISC
// @license.url https://www.isc.org/licenses/

// @host localhost:8002
// @BasePath /

// @tag.name Order
// @tag.description Requesting and placing orders

const dbURL = "postgresql://gobookstoreorders@localhost:5432/gobookstore"
const catalogURL = "localhost:8001"
const schema = "orders"
const bufsize = 1024 * 1024

var ctx = context.Background()

func main() {
	db, r := initSQL()
	defer db.Close()

	cConn, err := grpc.Dial(catalogURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to catalog service on %s: %s", catalogURL, err)
	}
	defer cConn.Close()

	lis, err := net.Listen("tcp", "localhost:8002")
	if err != nil {
		log.Fatalf("Failed to listen due to %s", err)
	}
	log.Println("Listening on localhost:8002")

	cc := catalog.NewCatalogClient(cConn)
	c := catalogservice.New(cc)
	s := orderservice.New(r, c)

	grpcServer := grpc.NewServer()
	orders.RegisterOrdersServer(grpcServer, server.New(s))
	log.Println("Starting gRPC server")
	grpcServer.Serve(lis)
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
