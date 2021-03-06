package grpc_test

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	authorInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/author/inmemory"
	authorservice "github.com/Vesninovich/go-tasks/book-store/catalog/author/service"
	bookInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/book/inmemory"
	bookservice "github.com/Vesninovich/go-tasks/book-store/catalog/book/service"
	categoryInMemory "github.com/Vesninovich/go-tasks/book-store/catalog/category/inmemory"
	categoryservice "github.com/Vesninovich/go-tasks/book-store/catalog/category/service"
	cataloggrpc "github.com/Vesninovich/go-tasks/book-store/catalog/grpc"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/catalog"
	pb "github.com/Vesninovich/go-tasks/book-store/common/catalog"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufsize = 1024 * 1024

var lis *bufconn.Listener
var ctx = context.Background()
var aut book.Author
var cats = []book.Category{
	{},
	{},
}
var as *authorservice.Service
var bs *bookservice.BookService
var cs *categoryservice.Service

func TestCreateBook(t *testing.T) {
	s := setup(t)
	defer s.GracefulStop()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %s", err)
	}
	defer conn.Close()
	client := pb.NewCatalogClient(conn)

	_, err = client.CreateBook(ctx, &catalog.BookCreateDTO{
		Name: "Test",
		Author: &catalog.Author{
			Id:   aut.ID[:],
			Name: aut.Name,
		},
		Categories: []*catalog.Category{
			{Id: cats[0].ID[:], Name: cats[0].Name},
			{Id: cats[1].ID[:], Name: cats[1].Name},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create valid book: %s", err)
	}

	_, err = client.CreateBook(ctx, &catalog.BookCreateDTO{
		Name: "Test",
		Author: &catalog.Author{
			Id:   aut.ID[:],
			Name: aut.Name,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create valid book with no categories: %s", err)
	}

	_, err = client.CreateBook(ctx, &catalog.BookCreateDTO{
		Name: "Test",
		Author: &catalog.Author{
			Id:   append(aut.ID[:], 0),
			Name: aut.Name,
		},
	})
	if err == nil {
		t.Error("Expected to get error for invalid author UUID")
	}

	_, err = client.CreateBook(ctx, &catalog.BookCreateDTO{
		Name: "Test",
		Author: &catalog.Author{
			Id:   aut.ID[:],
			Name: aut.Name,
		},
		Categories: []*catalog.Category{
			{Id: cats[0].ID[:], Name: cats[0].Name},
			{Id: append(cats[1].ID[:], 0), Name: cats[1].Name},
		},
	})
	if err == nil {
		t.Error("Expected to get error for invalid category UUID")
	}
}

func TestGetBooks(t *testing.T) {
	s := setup(t)
	defer s.GracefulStop()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %s", err)
	}
	defer conn.Close()
	client := pb.NewCatalogClient(conn)

	// Setup data
	a, err := as.CreateAuthor(ctx, "TestA2")
	if err != nil {
		t.Fatalf("Failed to create author: %s", err)
	}
	b0, err := client.CreateBook(ctx, &catalog.BookCreateDTO{
		Name: "TestA",
		Author: &catalog.Author{
			Id:   a.ID[:],
			Name: a.Name,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create valid book: %s", err)
	}
	b0ID, err := uuid.FromBytes(b0.Id)
	if err != nil {
		t.Fatalf("Failed to read uuid of book: %s", err)
	}
	b1, err := client.CreateBook(ctx, &catalog.BookCreateDTO{
		Name: "TestB",
		Author: &catalog.Author{
			Id:   aut.ID[:],
			Name: aut.Name,
		},
		Categories: []*catalog.Category{
			{Id: cats[0].ID[:], Name: cats[0].Name},
			{Id: cats[1].ID[:], Name: cats[1].Name},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create valid book: %s", err)
	}
	b1ID, err := uuid.FromBytes(b1.Id)
	if err != nil {
		t.Fatalf("Failed to read uuid of book: %s", err)
	}
	//

	stream, err := client.GetBooks(ctx, &pb.BooksQuery{})
	if err != nil {
		t.Fatalf("Failed to get books with empty query: %s", err)
	}
	count := 0
	for {
		b, err := stream.Recv()
		if err == io.EOF {
			if count != 2 {
				t.Fatalf("Wrong number of books read, expected 2, got %d", count)
			}
			break
		}
		if err != nil {
			t.Fatalf("Failed to read book from stream: %s", err)
		}
		_, err = uuid.FromBytes(b.Id)
		if err != nil {
			t.Fatalf("Failed to get uuid of book: %s", err)
		}
		count++
	}

	zero := uint32(0)
	one := uint32(1)
	two := uint32(2)
	stream, err = client.GetBooks(ctx, &pb.BooksQuery{
		From:   &zero,
		Count:  &one,
		Author: a.ID[:],
	})
	if err != nil {
		t.Fatalf("Failed to get books: %s", err)
	}
	count = 0
	for {
		b, err := stream.Recv()
		if err == io.EOF {
			if count != 1 {
				t.Fatalf("Wrong number of books read, expected 1, got %d", count)
			}
			break
		}
		if err != nil {
			t.Fatalf("Failed to read book from stream: %s", err)
		}
		id, err := uuid.FromBytes(b.Id)
		if err != nil {
			t.Fatalf("Failed to get uuid of book: %s", err)
		}
		if id != b0ID {
			t.Errorf("Got wrong book")
		}
		count++
	}

	stream, err = client.GetBooks(ctx, &pb.BooksQuery{
		From:       &zero,
		Count:      &one,
		Categories: [][]byte{cats[1].ID[:]},
	})
	if err != nil {
		t.Fatalf("Failed to get books: %s", err)
	}
	count = 0
	for {
		b, err := stream.Recv()
		if err == io.EOF {
			if count != 1 {
				t.Fatalf("Wrong number of books read, expected 1, got %d", count)
			}
			break
		}
		if err != nil {
			t.Fatalf("Failed to read book from stream: %s", err)
		}
		id, err := uuid.FromBytes(b.Id)
		if err != nil {
			t.Fatalf("Failed to get uuid of book: %s", err)
		}
		if id != b1ID {
			t.Errorf("Got wrong book")
		}
		count++
	}

	stream, err = client.GetBooks(ctx, &pb.BooksQuery{
		From:  &zero,
		Count: &two,
		Id:    b1.Id,
	})
	if err != nil {
		t.Fatalf("Failed to get books: %s", err)
	}
	count = 0
	for {
		b, err := stream.Recv()
		if err == io.EOF {
			if count != 1 {
				t.Fatalf("Wrong number of books read, expected 1, got %d", count)
			}
			break
		}
		if err != nil {
			t.Fatalf("Failed to read book from stream: %s", err)
		}
		id, err := uuid.FromBytes(b.Id)
		if err != nil {
			t.Fatalf("Failed to get uuid of book: %s", err)
		}
		if id != b1ID {
			t.Errorf("Got wrong book")
		}
		count++
	}
}

func setup(t *testing.T) *grpc.Server {
	lis = bufconn.Listen(bufsize)
	s := grpc.NewServer()

	as = authorservice.New(authorInMemory.New())
	cs = categoryservice.New(categoryInMemory.New())
	bs = bookservice.New(bookInMemory.New(), as, cs)

	var err error
	aut, err = as.CreateAuthor(ctx, "TestA")
	if err != nil {
		t.Fatalf("Failed to create author: %s", err)
	}
	cats[0], err = cs.CreateCategory(ctx, "TestC1", uuid.UUID{})
	if err != nil {
		t.Fatalf("Failed to create category: %s", err)
	}
	cats[1], err = cs.CreateCategory(ctx, "TestC2", uuid.UUID{})
	if err != nil {
		t.Fatalf("Failed to create category: %s", err)
	}

	pb.RegisterCatalogServer(s, cataloggrpc.New(bs))
	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %s", err)
		}
	}()

	return s
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
