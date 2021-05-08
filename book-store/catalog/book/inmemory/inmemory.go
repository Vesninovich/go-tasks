package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

// Repository represents in-memory repository of books
type Repository struct {
	authorRepo   author.Repository
	categoryRepo category.Repository

	data []book.Book
	lock sync.RWMutex
}

// New creates new in-memory repository of books
func New(authorRepo author.Repository, categoryRepo category.Repository) *Repository {
	return &Repository{
		authorRepo:   authorRepo,
		categoryRepo: categoryRepo,
		data:         make([]book.Book, 0),
	}
}

// GetAll gets all items from in-memory repository
func (r *Repository) GetAll(ctx context.Context) ([]book.Book, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.data[:], nil
}

// Get gets item by ID from in-memory repository
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (book.Book, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, item := range r.data {
		if item.ID == id {
			if !item.DeletedAt.IsZero() {
				return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
			}
			return item, nil
		}
	}

	return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
}

// Create creates item in in-memory repository
func (r *Repository) Create(ctx context.Context, dto bookrepo.CreateDTO) (book.Book, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	author, categories, err := r.storeChildren(ctx, dto.Author, dto.Categories)
	if err != nil {
		return book.Book{}, err
	}

	item := book.Book{
		ID:         uuid.New(),
		Name:       dto.Name,
		Author:     author,
		Categories: categories,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Time{},
		DeletedAt:  time.Time{},
	}
	r.data = append(r.data, item)
	return item, nil
}

// Update updates item in in-memory repository
func (r *Repository) Update(ctx context.Context, dto bookrepo.UpdateDTO) (book.Book, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, item := range r.data {
		if item.ID == dto.ID {
			if !item.DeletedAt.IsZero() {
				return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", dto.ID)}
			}
			author, categories, err := r.storeChildren(ctx, dto.Author, dto.Categories)
			if err != nil {
				return book.Book{}, err
			}
			r.data[i] = book.Book{
				ID:         dto.ID,
				Name:       dto.Name,
				Author:     author,
				Categories: categories,
				CreatedAt:  item.CreatedAt,
				UpdatedAt:  time.Now(),
				DeletedAt:  time.Time{},
			}
			return r.data[i], nil
		}
	}
	return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", dto.ID)}
}

// Delete deletes item in in-memory repository
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) (book.Book, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i, item := range r.data {
		if item.ID == id {
			if !item.DeletedAt.IsZero() {
				return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
			}
			r.data[i] = book.Book{
				ID:        id,
				Name:      item.Name,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: time.Now(),
			}
			return r.data[i], nil
		}
	}
	return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
}

func (r *Repository) storeChildren(
	ctx context.Context,
	authorDto author.UpdateDTO,
	categoriesDto []category.UpdateDTO,
) (author book.Author, categories []book.Category, err error) {
	author, err = r.storeAuthor(ctx, authorDto)
	if err != nil {
		return
	}
	categories, err = r.storeCategories(ctx, categoriesDto)
	return
}

func (r *Repository) storeAuthor(
	ctx context.Context,
	dto author.UpdateDTO,
) (aut book.Author, err error) {
	if dto.ID.IsZero() {
		aut, err = r.authorRepo.Create(ctx, author.CreateDTO{Name: dto.Name})
	} else {
		aut, err = r.authorRepo.Update(ctx, dto)
	}
	return
}

func (r *Repository) storeCategories(
	ctx context.Context,
	dtos []category.UpdateDTO,
) (categories []book.Category, err error) {
	categories = make([]book.Category, len(dtos))
	var cat book.Category
	for i, dto := range dtos {
		cat, err = r.storeCategory(ctx, dto)
		if err != nil {
			return
		}
		categories[i] = cat
	}
	return
}

func (r *Repository) storeCategory(
	ctx context.Context,
	dto category.UpdateDTO,
) (cat book.Category, err error) {
	if dto.ID.IsZero() {
		cat, err = r.categoryRepo.Create(ctx, category.CreateDTO{
			Name:     dto.Name,
			ParentID: dto.ParentID,
		})
	} else {
		cat, err = r.categoryRepo.Update(ctx, dto)
	}
	return
}
