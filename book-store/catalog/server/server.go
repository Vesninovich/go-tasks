package server

import (
	"context"

	bookservice "github.com/Vesninovich/go-tasks/book-store/catalog/book/service"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/catalog"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Server implements catalog gRPC server
type Server struct {
	catalog.UnimplementedCatalogServer

	bookService *bookservice.BookService
}

// New creates Server object
func New(bookService *bookservice.BookService) *Server {
	return &Server{
		bookService: bookService,
	}
}

// GetBooks streams books
func (s *Server) GetBooks(q *catalog.BooksQuery, stream catalog.Catalog_GetBooksServer) (err error) {
	bookID, autID, catIDs, err := getUUIDs(q.Id, q.Author, q.Categories)
	if err != nil {
		return
	}
	var from, count uint
	if q.From == nil {
		from = 0
	} else {
		from = uint(*q.From)
	}
	if q.Count == nil {
		count = 0
	} else {
		count = uint(*q.Count)
	}
	data, err := s.bookService.GetBooks(context.Background(), from, count, book.Query{
		ID:         bookID,
		Author:     autID,
		Categories: catIDs,
	})
	for _, item := range data {
		err = stream.Send(makeBookResponse(item))
		if err != nil {
			return
		}
	}
	return
}

// CreateBook creates book
func (s *Server) CreateBook(ctx context.Context, dto *catalog.BookCreateDTO) (*catalog.Book, error) {
	_, autID, catIDs, err := getUUIDs(nil, dto.Author, dto.Categories)
	if err != nil {
		return nil, err
	}
	cats := make([]book.Category, len(catIDs))
	for i, id := range catIDs {
		cats[i] = book.Category{ID: id}
	}
	b, err := s.bookService.CreateBook(ctx, dto.Name, book.Author{ID: autID}, cats)
	if err != nil {
		return nil, err
	}
	return makeBookResponse(b), err
}

func getUUIDs(bID []byte, author []byte, categories [][]byte) (bookID uuid.UUID, autID uuid.UUID, catIDs []uuid.UUID, err error) {
	if bID != nil && len(bID) != 0 {
		bookID, err = uuid.FromBytes(bID)
		if err != nil {
			return
		}
	}
	if author != nil && len(author) != 0 {
		autID, err = uuid.FromBytes(author)
		if err != nil {
			return
		}
	}
	var id uuid.UUID
	catIDs = make([]uuid.UUID, len(categories))
	for i, c := range categories {
		id, err = uuid.FromBytes(c)
		if err != nil {
			return
		}
		catIDs[i] = id
	}
	return
}

func makeBookResponse(item book.Book) *catalog.Book {
	categories := make([]*catalog.Category, len(item.Categories))
	for i, cat := range item.Categories {
		categories[i] = &catalog.Category{
			Id:       cat.ID[:],
			Name:     cat.Name,
			ParentId: cat.ParentID[:],
		}
	}

	return &catalog.Book{
		Id:   item.ID[:],
		Name: item.Name,
		Author: &catalog.Author{
			Id:   item.Author.ID[:],
			Name: item.Author.Name,
		},
		Categories: categories,
	}
}
