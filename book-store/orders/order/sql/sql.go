package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	"github.com/Vesninovich/go-tasks/book-store/orders/order"
	"github.com/jmoiron/sqlx"
)

// Repository provides access to relational DB storage of orders
type Repository struct {
	db     *sqlx.DB
	schema string
}

type fromDB struct {
	ID          string
	Description string
	BookID      string `db:"book_id"`
}

// New creates a new instance of Repository
func New(db *sqlx.DB, schema string) *Repository {
	return &Repository{db, schema}
}

// CreateTableStmt of orders
func (r *Repository) CreateTableStmt() string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.orders(
  id uuid PRIMARY KEY,
  description text NOT NULL,
	book_id uuid,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);`, r.schema)
}

// GetAll gets all non-deleted orders
func (r *Repository) GetAll(ctx context.Context) (orders []order.DTO, err error) {
	data := []fromDB{}
	err = r.db.SelectContext(ctx, &data, fmt.Sprintf("SELECT id, description, book_id FROM %s.orders WHERE deleted_at=$1;", r.schema), time.Time{})
	if err != nil {
		return
	}
	var dto order.DTO
	orders = make([]order.DTO, len(data))
	for i, item := range data {
		dto, err = item.toDTO()
		if err != nil {
			return
		}
		orders[i] = dto
	}
	return
}

// Get gets non-deleted order by ID
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (order.DTO, error) {
	o := fromDB{}
	err := r.db.GetContext(
		ctx, &o, fmt.Sprintf("SELECT id, description, book_id FROM %s.orders WHERE id=$1 AND deleted_at=$2;", r.schema), id.String(), time.Time{},
	)
	if err == sql.ErrNoRows {
		return order.DTO{}, &commonerrors.NotFound{What: fmt.Sprintf("Order with ID %s", id)}
	}
	if err != nil {
		return order.DTO{}, nil
	}
	return o.toDTO()
}

// Create stores new order
func (r *Repository) Create(ctx context.Context, dto order.CreateDTO) (order.DTO, error) {
	id := uuid.New()
	_, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf(`INSERT INTO %s.orders (id, description, book_id, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6)`, r.schema),
		id.String(), dto.Description, dto.BookID.String(), time.Now(), time.Time{}, time.Time{},
	)
	return order.DTO{
		ID:        id,
		CreateDTO: dto,
	}, err
}

// Update updates stored non-deleted order
func (r *Repository) Update(ctx context.Context, dto order.DTO) (order.DTO, error) {
	res, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf(`UPDATE %s.orders
			SET description=$3, book_id=$4, updated_at=$5
			WHERE id=$1 AND deleted_at=$2;`, r.schema),
		dto.ID.String(), time.Time{}, dto.Description, dto.BookID.String(), time.Now(),
	)
	if err != nil {
		return order.DTO{}, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return order.DTO{}, err
	}
	if count == 0 {
		return order.DTO{}, &commonerrors.NotFound{What: fmt.Sprintf("Order with ID %s", dto.ID)}
	}
	return dto, err
}

// Delete sets stored order with id as deleted
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) (order.DTO, error) {
	var o fromDB
	err := r.db.GetContext(
		ctx,
		&o,
		fmt.Sprintf(`SELECT id, description, book_id FROM %s.orders WHERE id=$1 AND deleted_at=$2;`, r.schema),
		id.String(), time.Time{},
	)
	if err == sql.ErrNoRows {
		return order.DTO{}, &commonerrors.NotFound{What: fmt.Sprintf("Order with ID %s", id)}
	}
	if err != nil {
		return order.DTO{}, err
	}
	res, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf(`UPDATE %s.orders
			SET deleted_at=$2
			WHERE id=$1 AND deleted_at=$3;`, r.schema),
		id.String(), time.Now(), time.Time{},
	)
	if err != nil {
		return order.DTO{}, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return order.DTO{}, err
	}
	if count == 0 {
		return order.DTO{}, &commonerrors.NotFound{What: fmt.Sprintf("Order with ID %s", id)}
	}
	return o.toDTO()
}

func (f fromDB) toDTO() (order.DTO, error) {
	id, err := uuid.FromString(f.ID)
	if err != nil {
		return order.DTO{}, err
	}
	bID, err := uuid.FromString(f.BookID)
	return order.DTO{
		ID: id,
		CreateDTO: order.CreateDTO{
			Description: f.Description,
			BookID:      bID,
		},
	}, err
}
