# TODO

## Запуск

`docker-compose up`

## todo

Дописать тесты

## Задание

Создать ToDo приложение с использованием REST API. Для хранения данных использовать любую реляционную базу данных (MySQL, PostgreSQL).
Главная сущность приложения `Task` содержит в себе следующие поля:
- `Name` - название
- `Description` - описание
- `DueDate` - дедлайн
- `Status` - статус (новый, выполняется, отменен, выполнен, просрочен)

Необходимо разработать 5 endpoints:
- Для создания одного таска
- Для чтения данных одного таска
- Для обновления данных одного таска
- Для удаления таска
- Для считывания нескольких тасков

Рекомендуется использовать стандартные библиотеки golang для работы с REST и SQL. Формат обмена данными - JSON.
