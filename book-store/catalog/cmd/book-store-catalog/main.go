package main

import (
	"fmt"
	"log"
	"net"

	authorservice "github.com/Vesninovich/go-tasks/book-store/catalog/author/service"
	authorsql "github.com/Vesninovich/go-tasks/book-store/catalog/author/sql"
	bookservice "github.com/Vesninovich/go-tasks/book-store/catalog/book/service"
	booksql "github.com/Vesninovich/go-tasks/book-store/catalog/book/sql"
	categoryservice "github.com/Vesninovich/go-tasks/book-store/catalog/category/service"
	categorysql "github.com/Vesninovich/go-tasks/book-store/catalog/category/sql"
	cataloggrpc "github.com/Vesninovich/go-tasks/book-store/catalog/grpc"
	pb "github.com/Vesninovich/go-tasks/book-store/common/catalog"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

// @title Book Store Catalog Service
// @version 0.0
// @description Service for creating and quering books catalog

// @contact.name Dimas
// @contact.url https://github.com/Vesninovich
// @contact.email dmitry@vesnin.work

// @license.name ISC
// @license.url https://www.isc.org/licenses/

// @host localhost:8001
// @BasePath /

// @tag.name Book
// @tag.description Quering and creating books

const dbURL = "postgresql://gobookstorecatalog@localhost:5432/gobookstore"
const schema = "catalog"

func main() {
	db, ar, cr, br := initSQL()
	defer db.Close()

	lis, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatalf("Failed to listen due to %s", err)
	}
	log.Println("Listening on localhost:8001")

	grpcServer := grpc.NewServer()

	as := authorservice.New(ar)
	cs := categoryservice.New(cr)
	bs := bookservice.New(br, as, cs)

	pb.RegisterCatalogServer(grpcServer, cataloggrpc.New(bs))
	log.Println("Starting gRPC server")
	grpcServer.Serve(lis)
}

func initSQL() (*sqlx.DB, *authorsql.Repository, *categorysql.Repository, *booksql.Repository) {
	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB at URL %s\n%s", dbURL, err)
	}
	log.Printf("Connected to DB at URL %s\n", dbURL)

	log.Println("Creating schema")
	s := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schema)
	log.Println(s)
	db.MustExec(s)

	a := authorsql.New(db, schema)
	c := categorysql.New(db, schema)
	b := booksql.New(db, schema)

	log.Println("Creating tables")
	log.Println(a.CreateTableStmt())
	db.MustExec(a.CreateTableStmt())
	log.Println(c.CreateTableStmt())
	db.MustExec(c.CreateTableStmt())
	log.Println(b.CreateTableStmt())
	db.MustExec(b.CreateTableStmt())
	log.Println("Finished setting up DB")

	return db, a, c, b
}
