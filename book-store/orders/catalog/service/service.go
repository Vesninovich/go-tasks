package service

import (
	"context"
	"fmt"
	"io"

	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/catalog"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Service for accessing the catalog
type Service struct {
	catClient catalog.CatalogClient
	cache     []book.Book // used for in-memory testing, may be used later
}

// New creates service
func New(c catalog.CatalogClient) *Service {
	return &Service{catClient: c, cache: make([]book.Book, 0)}
}

// GetBook gets book from catalog
func (s *Service) GetBook(ctx context.Context, id uuid.UUID) (bk book.Book, err error) {
	cl, err := s.catClient.GetBooks(ctx, &catalog.BooksQuery{
		Id: id[:],
	})
	if err != nil {
		return
	}
	res, err := cl.Recv()
	if err == io.EOF {
		return bk, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
	}
	if err != nil {
		return
	}
	bk, err = resToBook(res)
	return
}

func resToBook(res *catalog.Book) (bk book.Book, err error) {
	id, err := uuid.FromBytes(res.Id)
	if err != nil {
		return
	}
	aID, err := uuid.FromBytes(res.Author.Id)
	if err != nil {
		return
	}
	cats := make([]book.Category, len(res.Categories))
	var cID, cPID uuid.UUID
	for i, c := range res.Categories {
		cID, err = uuid.FromBytes(c.Id)
		if err != nil {
			return
		}
		cPID, err = uuid.FromBytes(c.ParentId)
		if err != nil {
			return
		}
		cats[i] = book.Category{
			ID:       cID,
			Name:     c.Name,
			ParentID: cPID,
		}
	}
	bk = book.Book{
		ID:   id,
		Name: res.Name,
		Author: book.Author{
			ID:   aID,
			Name: res.Author.Name,
		},
		Categories: cats,
	}
	return
}

// GetBookFromCache gets book from cache
func (s *Service) getBookFromCache(ctx context.Context, id uuid.UUID) (book.Book, error) {
	for _, b := range s.cache {
		if b.ID == id {
			return b, nil
		}
	}
	return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
}

// CheckBookInCache check if book in cache
func (s *Service) checkBookInCache(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := s.getBookFromCache(ctx, id)
	if err == nil {
		return true, err
	}
	if _, notFound := err.(*commonerrors.NotFound); notFound {
		return false, nil
	}
	return false, err
}
