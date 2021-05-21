module github.com/Vesninovich/go-tasks/book-store/catalog

go 1.16

require (
	github.com/Vesninovich/go-tasks/book-store/common v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/jmoiron/sqlx v1.3.4
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.37.1
)

replace github.com/Vesninovich/go-tasks/book-store/common => ../common
