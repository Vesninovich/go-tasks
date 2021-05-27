# Book Store Orders

## Testing

`go test ./...` прогнать юнит-тесты

`go test ./... -tags=sql` прогнать тесты с базой

`go test ./... -tags=integr_full` прогнать полный интеграционный тест (нужно, что сервис каталога был запущен)
