package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	bookrepo "github.com/Vesninovich/go-tasks/book-store/catalog/book"
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	"github.com/jmoiron/sqlx"
)

// Table of authors
const Table = `
CREATE TABLE IF NOT EXISTS books(
  id uuid PRIMARY KEY,
  name text NOT NULL,
	author_id uuid REFERENCES authors,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS books_categories(
	book_id uuid REFERENCES books ON DELETE CASCADE,
	category_id uuid REFERENCES categories ON DELETE CASCADE,
	PRIMARY KEY (book_id, category_id)
);`

// Repository provides access to relational DB storage of tasks
type Repository struct {
	db *sqlx.DB
}

type fromDB struct {
	ID         string
	Name       string
	AuthorID   string `db:"author_id"`
	AuthorName string `db:"author_name"`
}

type catsFromDB struct {
	ID               string
	CategoryID       string         `db:"category_id"`
	CategoryName     string         `db:"category_name"`
	CategoryParentID sql.NullString `db:"category_parent_id"`
}

// New creates a new instance of SQLRepository
func New(db *sqlx.DB) *Repository {
	return &Repository{db}
}

// Get gets
func (r *Repository) Get(ctx context.Context, from, count uint, query book.Query) ([]book.Book, error) {
	data := []fromDB{}
	err := r.db.SelectContext(
		ctx, &data, getSelectBooksStatement(from, count, query), time.Time{},
	)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []book.Book{}, err
	}
	catData := []catsFromDB{}
	err = r.db.SelectContext(
		ctx, &catData, getSelectCategoriesStatement(data), time.Time{},
	)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return mapBookData(data, catData)
}

// TODO: rewrite to string builder or squirrel
func getSelectBooksStatement(from, count uint, query book.Query) (stmt string) {
	stmt = `SELECT
		books.id as id,
		books.name as name,
		authors.id as author_id,
		authors.name as author_name
		FROM books `
	stmt += `
		INNER JOIN authors
		ON (authors.id=books.author_id)`
	qStart := `
		WHERE `
	qCont := `
		AND `
	if !query.ID.IsZero() {
		stmt += fmt.Sprintf(`%sbooks.id='%s'`, qStart, query.ID)
		qStart = qCont
	}
	if !query.Author.IsZero() {
		stmt += fmt.Sprintf(`%sbooks.author_id='%s'`, qStart, query.Author)
		qStart = qCont
	}
	if len(query.Categories) != 0 {
		// is this portable?
		stmt += fmt.Sprintf(`%sEXISTS (
			SELECT 1
			FROM books_categories as bc
			WHERE bc.book_id=books.id
			AND bc.category_id IN (
				'%s'`, qStart, query.Categories[0])
		for _, c := range query.Categories[1:] {
			stmt += fmt.Sprintf(`,
				'%s'`, c)
		}
		stmt += `
			)
		)`
		qStart = qCont
	}
	stmt += qStart + "books.deleted_at=$1 AND authors.deleted_at=$1"
	if from != 0 && query.ID.IsZero() {
		stmt += fmt.Sprintf(`
		OFFSET %d`, from)
	}
	if !query.ID.IsZero() {
		count = 1
	}
	if count != 0 {
		stmt += fmt.Sprintf(`
		LIMIT %d`, count)
	}
	log.Println(stmt)
	return
}

func getSelectCategoriesStatement(bookData []fromDB) (stmt string) {
	stmt = fmt.Sprintf(`SELECT
		books.id as id,
		categories.id as category_id,
		categories.name as category_name,
		categories.parent_id as category_parent_id
		FROM books
		INNER JOIN categories
		ON EXISTS (
	 		SELECT 1
	 		FROM books_categories as bc
		 	WHERE bc.book_id=books.id AND bc.category_id=category_id
		)
		WHERE categories.deleted_at=$1
		AND books.id IN (
			'%s'`, bookData[0].ID)
	// is this portable?
	for _, b := range bookData[1:] {
		stmt += fmt.Sprintf(`,
			'%s'`, b.ID)
	}
	stmt += `
		);`
	log.Println(stmt)
	return
}

func mapBookData(booksData []fromDB, catData []catsFromDB) (books []book.Book, err error) {
	booksMap := make(map[string]*book.Book)
	var aID, bID, cID, cPID uuid.UUID
	var bk *book.Book
	var exists bool
	for _, b := range booksData {
		bID, err = uuid.FromString(b.ID)
		if err != nil {
			return
		}
		aID, err = uuid.FromString(b.AuthorID)
		if err != nil {
			return
		}
		bk = &book.Book{
			ID:   bID,
			Name: b.Name,
			Author: book.Author{
				ID:   aID,
				Name: b.AuthorName,
			},
			Categories: make([]book.Category, 0),
		}
		booksMap[b.ID] = bk
	}
	for _, c := range catData {
		bk, exists = booksMap[c.ID]
		if !exists {
			err = &commonerrors.NotFound{What: "Some strange shit"}
			return
		}
		cID, err = uuid.FromString(c.CategoryID)
		if err != nil {
			return
		}
		cPID = uuid.UUID{}
		if c.CategoryParentID.Valid {
			cPID, err = uuid.FromString(c.CategoryParentID.String)
			if err != nil {
				return
			}
		}
		bk.Categories = append(bk.Categories, book.Category{
			ID:       cID,
			Name:     c.CategoryName,
			ParentID: cPID,
		})
	}
	books = make([]book.Book, 0, len(booksData))
	for _, b := range booksMap {
		books = append(books, *b)
	}
	return
}

// Create creates
func (r *Repository) Create(ctx context.Context, dto bookrepo.CreateDTO) (book.Book, error) {
	id := uuid.New()
	idStr := id.String()
	tx, err := r.db.BeginTxx(ctx, nil) // TODO: correct options
	if err != nil {
		return book.Book{}, err
	}
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO books (id, name, author_id, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
		idStr, dto.Name, dto.Author.ID.String(), time.Now(), time.Time{}, time.Time{},
	)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			err = rbErr
		}
		return book.Book{}, err
	}
	for _, cat := range dto.Categories {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO books_categories (book_id, category_id)
				VALUES ($1, $2)`,
			idStr, cat.ID.String(),
		)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = rbErr
			}
			return book.Book{}, err
		}
	}
	err = tx.Commit()
	return book.Book{ID: id, Name: dto.Name, Author: dto.Author, Categories: dto.Categories}, err
}

// Update updates
func (r *Repository) Update(ctx context.Context, dto book.Book) (b book.Book, err error) {
	books, err := r.Get(ctx, 0, 1, book.Query{ID: dto.ID})
	if err != nil {
		return
	}
	if len(books) == 0 {
		return b, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", dto.ID)}
	}
	b = books[0]
	catsEq := equalCategories(dto, b)
	if dto.Name == b.Name && dto.Author.ID == b.Author.ID && catsEq {
		return
	}
	idStr := dto.ID.String()
	tx, err := r.db.BeginTxx(ctx, nil) // TODO: correct options
	if err != nil {
		return
	}
	_, err = tx.ExecContext(
		ctx,
		`UPDATE books
			SET name=$2, author_id=$3, updated_at=$4
			WHERE id=$1`,
		idStr, dto.Name, dto.Author.ID.String(), time.Now(),
	)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			err = rbErr
		}
		return
	}
	if !catsEq {
		_, err = tx.ExecContext(
			ctx,
			`DELETE FROM books_categories
				WHERE book_id=$1`,
			idStr,
		)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = rbErr
			}
			return
		}
		for _, cat := range dto.Categories {
			_, err := tx.ExecContext(
				ctx,
				`INSERT INTO books_categories (book_id, category_id)
					VALUES ($1, $2)`,
				idStr, cat.ID.String(),
			)
			if err != nil {
				rbErr := tx.Rollback()
				if rbErr != nil {
					err = rbErr
				}
				return book.Book{}, err
			}
		}
	}
	b = dto
	err = tx.Commit()
	return
}

func equalCategories(a book.Book, b book.Book) bool {
	if len(a.Categories) != len(b.Categories) {
		return false
	}
	for i, c := range a.Categories {
		if c.ID != b.Categories[i].ID {
			return false
		}
	}
	return true
}

// Delete deletes
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) (book.Book, error) {
	books, err := r.Get(ctx, 0, 1, book.Query{ID: id})
	if err != nil {
		return book.Book{}, err
	}
	if len(books) == 0 {
		return book.Book{}, &commonerrors.NotFound{What: fmt.Sprintf("Book with ID %s", id)}
	}
	idStr := id.String()
	tx, err := r.db.BeginTxx(ctx, nil) // TODO: correct options
	if err != nil {
		return book.Book{}, err
	}
	_, err = tx.ExecContext(
		ctx,
		`UPDATE books
			SET deleted_at=$2
			WHERE id=$1`,
		idStr, time.Now(),
	)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			err = rbErr
		}
		return book.Book{}, err
	}
	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM books_categories
			WHERE book_id=$1`,
		idStr,
	)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			err = rbErr
		}
		return book.Book{}, err
	}
	err = tx.Commit()
	return books[0], err
}
