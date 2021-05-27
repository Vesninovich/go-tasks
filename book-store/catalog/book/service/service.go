package bookservice

import (
	"context"

	authorservice "github.com/Vesninovich/go-tasks/book-store/catalog/author/service"
	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	categoryservice "github.com/Vesninovich/go-tasks/book-store/catalog/category/service"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
)

// BookService is service for interacting with books
type BookService struct {
	bookRepo        bookrepo.Repository
	authorService   *authorservice.Service
	categoryService *categoryservice.Service
}

// New creates new BookService
func New(bookRepo bookrepo.Repository, authorService *authorservice.Service, categoryService *categoryservice.Service) *BookService {
	return &BookService{
		bookRepo:        bookRepo,
		authorService:   authorService,
		categoryService: categoryService,
	}
}

// GetBooks fetches count saved books from some number according to query
func (s *BookService) GetBooks(ctx context.Context, from, count uint, query book.Query) ([]book.Book, error) {
	if count == 0 {
		count = 10
	}
	return s.bookRepo.Get(ctx, from, count, query)
}

// CreateBook saves new book if name is not empty, listed author and all categories exist
func (s *BookService) CreateBook(ctx context.Context, name string, aut book.Author, cats []book.Category) (book.Book, error) {
	if name == "" {
		return book.Book{}, &commonerrors.InvalidInput{Reason: "name is required"}
	}
	if aut.ID.IsZero() {
		// TODO: optimize nested creation
		a, err := s.authorService.CreateAuthor(ctx, aut.Name)
		if err != nil {
			return book.Book{}, err
		}
		aut = a
	} else {
		_, err := s.authorService.GetAuthor(ctx, aut.ID)
		if err != nil {
			return book.Book{}, err
		}
	}
	// TODO: optimize check and nested creation
	for i, cat := range cats {
		if cat.ID.IsZero() {
			c, err := s.categoryService.CreateCategory(ctx, cat.Name, cat.ParentID)
			if err != nil {
				return book.Book{}, err
			}
			cats[i] = c
		} else {
			_, err := s.categoryService.GetCategory(ctx, cat.ID)
			if err != nil {
				return book.Book{}, err
			}
		}
	}
	return s.bookRepo.Create(ctx, bookrepo.CreateDTO{
		Name:       name,
		Author:     aut,
		Categories: cats,
	})
}
