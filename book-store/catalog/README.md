# Book Store Catalog

## Running

`go run ./cmd/book-store-catalog/main.go`

Надо, чтобы постгрес крутился и слушал на 5432 и чтобы был свободен порт 8001

## Testing

`go test ./...` прогнать юнит-тесты

`go test ./... -tags=sql` прогнать тесты с базой
