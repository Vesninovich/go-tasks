# Book Store Catalog

## Running

`go run ./cmd/book-store-catalog/main.go`

Надо, чтобы постгрес крутился и слушал на 5432 и чтобы был свободен порты 8001-2

## [Swagger](http://localhost:8002/book/swagger)

## Testing

`go test ./...` прогнать юнит-тесты

`go test ./... -tags=sql` прогнать тесты с базой

## TODO

- REST API на авторов и категории
- тесты REST API
- конфигурация
