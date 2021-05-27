module github.com/Vesninovich/go-tasks/book-store/orders

go 1.16

require (
	github.com/Vesninovich/go-tasks/book-store/common v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v4 v4.11.0
	github.com/jmoiron/sqlx v1.3.4
	google.golang.org/grpc v1.38.0 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace github.com/Vesninovich/go-tasks/book-store/common => ../common
