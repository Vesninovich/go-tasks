# Книжный магазин

2 микросервиса с API Gateway.

1. Catalog service - хранит сущности book (uuid, name, author_id), author(uuid, name), category(uuid, name, parent_uuid).
Book to category many-to-many. У списка books будут фильтры по author и category, pagination

2. Order service - хранит сущность order(uuid, book_id, description)

У всех сущностей есть (created_at, updated_at, deleted_at)

Между собой микросервисы общаются на gRPC. Для каждого микросервиса свой API Gateway c REST API и сваггером. У каждого микросервиса своя БД. База данных MySQL или PostgreSQL. Для работы можно использовать sqlx, вместо стандартной библиотеки sql. Можно также использовать Squirrel в качестве для построения SQL запросов.

Полезные ссылки:
- https://grpc.io/
- https://grpc.io/docs/languages/go/ 
- https://github.com/Masterminds/squirrel
- https://github.com/jmoiron/sqlx
- https://github.com/swaggo/swag 
