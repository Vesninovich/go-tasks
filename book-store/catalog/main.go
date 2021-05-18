package main

import (
	"log"
	"net"

	authorInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/author/inmemory"
	authorservice "github.com/Vesninovich/go-tasks/book-store/catalog/author/service"
	bookInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/book/inmemory"
	bookservice "github.com/Vesninovich/go-tasks/book-store/catalog/book/service"
	categoryInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/category/inmemory"
	categoryservice "github.com/Vesninovich/go-tasks/book-store/catalog/category/service"
	"github.com/Vesninovich/go-tasks/book-store/catalog/server"
	pb "github.com/Vesninovich/go-tasks/book-store/common/catalog"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:1488")
	if err != nil {
		log.Fatalf("Failed to listen due to %s", err)
	}
	log.Println("Listening on localhost:1488")

	grpcServer := grpc.NewServer()

	as := authorservice.New(authorInMemory.New())
	cs := categoryservice.New(categoryInMemory.New())
	bs := bookservice.New(bookInMemory.New(), as, cs)

	pb.RegisterCatalogServer(grpcServer, server.New(bs))
	log.Println("Starting gRPC server")
	grpcServer.Serve(lis)
}
