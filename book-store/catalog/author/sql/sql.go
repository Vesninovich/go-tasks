package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Vesninovich/go-tasks/book-store/catalog/author"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	"github.com/jmoiron/sqlx"
)

// Table of authors
const Table = `
CREATE TABLE IF NOT EXISTS authors(
  id uuid PRIMARY KEY,
  name text NOT NULL,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);`

// Repository provides access to relational DB storage of tasks
type Repository struct {
	db *sqlx.DB
}

type fromDB struct {
	ID   string
	Name string
}

// New creates a new instance of SQLRepository
func New(db *sqlx.DB) *Repository {
	return &Repository{db}
}

// GetAll gets all non-deleted authors
func (r *Repository) GetAll(ctx context.Context) (authors []book.Author, err error) {
	data := []fromDB{}
	err = r.db.SelectContext(ctx, &data, "SELECT id, name FROM authors WHERE deleted_at=$1;", time.Time{})
	if err != nil {
		return
	}
	var id uuid.UUID
	authors = make([]book.Author, len(data))
	for i, item := range data {
		id, err = uuid.FromString(item.ID)
		if err != nil {
			return
		}
		authors[i] = book.Author{
			ID:   id,
			Name: item.Name,
		}
	}
	return
}

// Get gets non-deleted author by ID
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (book.Author, error) {
	a := fromDB{}
	err := r.db.GetContext(
		ctx, &a, "SELECT id, name FROM authors WHERE id=$1 AND deleted_at=$2;", id.String(), time.Time{},
	)
	if err == sql.ErrNoRows {
		return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
	}
	if err != nil {
		return book.Author{}, nil
	}
	foundID, err := uuid.FromString(a.ID)
	return book.Author{ID: foundID, Name: a.Name}, err
}

// Create stores new author
func (r *Repository) Create(ctx context.Context, dto author.CreateDTO) (book.Author, error) {
	id := uuid.New()
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO authors (id, name, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5)`,
		id.String(), dto.Name, time.Now(), time.Time{}, time.Time{},
	)
	return book.Author{ID: id, Name: dto.Name}, err
}

// Update updates stored non-deleted author
func (r *Repository) Update(ctx context.Context, dto book.Author) (book.Author, error) {
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE authors
			SET name=$3, updated_at=$4
			WHERE id=$1 AND deleted_at=$2;`,
		dto.ID.String(), time.Time{}, dto.Name, time.Now(),
	)
	if err != nil {
		return book.Author{}, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return book.Author{}, err
	}
	if count == 0 {
		return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", dto.ID)}
	}
	return dto, err
}

// Delete sets stored author with id as deleted
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) (book.Author, error) {
	var a fromDB
	err := r.db.QueryRowxContext(
		ctx,
		`SELECT id, name FROM authors WHERE id=$1 AND deleted_at=$2;`,
		id.String(), time.Time{},
	).Scan(&a.ID, &a.Name)
	if err == sql.ErrNoRows {
		return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
	}
	if err != nil {
		return book.Author{}, err
	}
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE authors
			SET deleted_at=$2
			WHERE id=$1 AND deleted_at=$3;`,
		id.String(), time.Now(), time.Time{},
	)
	if err != nil {
		return book.Author{}, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return book.Author{}, err
	}
	if count == 0 {
		return book.Author{}, &commonerrors.NotFound{What: fmt.Sprintf("Author with ID %s", id)}
	}
	return book.Author{ID: id, Name: a.Name}, err
}
