# Book Store Orders

## Running

`go run ./cmd/book-store-orders/main.go`

Надо, чтобы постгрес крутился и слушал на 5432, чтобы grpc каталога крутился и слушал на 8001 и чтобы были свободны порты 8003-4

## [Swagger](http://localhost:8004/order/swagger)

## Testing

`go test ./...` прогнать юнит-тесты

`go test ./... -tags=sql` прогнать тесты с базой

`go test ./... -tags=integr_full` прогнать полный интеграционный тест (нужно, что сервис каталога был запущен)

## TODO

- тесты
- вынести сваггер сервер отдельно
- конфигурация
