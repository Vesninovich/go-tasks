package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Vesninovich/go-tasks/book-store/catalog/category"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	"github.com/jmoiron/sqlx"
)

// Table of categories
const Table = `
CREATE TABLE IF NOT EXISTS categories(
  id uuid PRIMARY KEY,
  name text NOT NULL,
	parent_id uuid REFERENCES categories DEFAULT NULL,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);`

// Repository provides access to relational DB storage of tasks
type Repository struct {
	db *sqlx.DB
}

type fromDB struct {
	ID       string
	Name     string
	ParentID sql.NullString `db:"parent_id"`
}

// New creates a new instance of SQLRepository
func New(db *sqlx.DB) *Repository {
	return &Repository{db}
}

// GetAll gets all non-deleted categories
func (r *Repository) GetAll(ctx context.Context) (categories []book.Category, err error) {
	data := []fromDB{}
	err = r.db.SelectContext(ctx, &data, "SELECT id, name, parent_id FROM categories WHERE deleted_at=$1;", time.Time{})
	if err != nil {
		return
	}
	var id, parentID uuid.UUID
	categories = make([]book.Category, len(data))
	for i, item := range data {
		id, err = uuid.FromString(item.ID)
		if err != nil {
			return
		}
		if item.ParentID.Valid {
			parentID, err = uuid.FromString(item.ParentID.String)
			if err != nil {
				return
			}
		}
		categories[i] = book.Category{
			ID:       id,
			Name:     item.Name,
			ParentID: parentID,
		}
	}
	return
}

// Get gets non-deleted category by ID
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (book.Category, error) {
	a := fromDB{}
	err := r.db.GetContext(
		ctx, &a, "SELECT id, name, parent_id FROM categories WHERE id=$1 AND deleted_at=$2;", id.String(), time.Time{},
	)
	if err == sql.ErrNoRows {
		return book.Category{}, &commonerrors.NotFound{What: fmt.Sprintf("Category with ID %s", id)}
	}
	if err != nil {
		return book.Category{}, nil
	}
	foundID, err := uuid.FromString(a.ID)
	if err != nil {
		return book.Category{}, err
	}
	parentID, err := uuid.FromString(a.ID)
	return book.Category{ID: foundID, Name: a.Name, ParentID: parentID}, err
}

// Create stores new category
func (r *Repository) Create(ctx context.Context, dto category.CreateDTO) (book.Category, error) {
	id := uuid.New()
	var parentID sql.NullString
	if !dto.ParentID.IsZero() {
		parentID.String = dto.ParentID.String()
		parentID.Valid = true
	}
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO categories (id, name, parent_id, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
		id.String(), dto.Name, parentID, time.Now(), time.Time{}, time.Time{},
	)
	return book.Category{ID: id, Name: dto.Name, ParentID: dto.ParentID}, err
}

// Update updates stored non-deleted category
func (r *Repository) Update(ctx context.Context, dto book.Category) (book.Category, error) {
	var parentID sql.NullString
	if !dto.ParentID.IsZero() {
		parentID.String = dto.ParentID.String()
		parentID.Valid = true
	}
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE categories
			SET name=$4, parent_id=$5, updated_at=$3
			WHERE id=$1 AND deleted_at=$2;`,
		dto.ID.String(), time.Time{}, time.Now(), dto.Name, parentID,
	)
	if err != nil {
		return book.Category{}, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return book.Category{}, err
	}
	if count == 0 {
		return book.Category{}, &commonerrors.NotFound{What: fmt.Sprintf("Category with ID %s", dto.ID)}
	}
	return dto, err
}

// Delete sets stored category with id as deleted
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) (book.Category, error) {
	var a fromDB
	err := r.db.QueryRowxContext(
		ctx,
		`SELECT id, name, parent_id FROM categories WHERE id=$1 AND deleted_at=$2;`,
		id.String(), time.Time{},
	).Scan(&a.ID, &a.Name, &a.ParentID)
	if err == sql.ErrNoRows {
		return book.Category{}, &commonerrors.NotFound{What: fmt.Sprintf("Category with ID %s", id)}
	}
	if err != nil {
		return book.Category{}, err
	}
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE categories
			SET deleted_at=$2
			WHERE id=$1 AND deleted_at=$3;`,
		id.String(), time.Now(), time.Time{},
	)
	if err != nil {
		return book.Category{}, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return book.Category{}, err
	}
	if count == 0 {
		return book.Category{}, &commonerrors.NotFound{What: fmt.Sprintf("Category with ID %s", id)}
	}
	var parentID uuid.UUID
	if a.ParentID.Valid {
		parentID, err = uuid.FromString(a.ParentID.String)
	}
	return book.Category{ID: id, Name: a.Name, ParentID: parentID}, err
}
