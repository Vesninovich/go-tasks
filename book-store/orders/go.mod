module github.com/Vesninovich/go-tasks/book-store/orders

go 1.16

require (
	github.com/Vesninovich/go-tasks/book-store/common v0.0.0-00010101000000-000000000000
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/swaggo/swag v1.7.0
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	golang.org/x/tools v0.1.2 // indirect
	google.golang.org/grpc v1.38.0
)

replace github.com/Vesninovich/go-tasks/book-store/common => ../common
