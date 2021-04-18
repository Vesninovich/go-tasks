package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/Vesninovich/go-tasks/todos/httpserver"
	taskhttp "github.com/Vesninovich/go-tasks/todos/task/http"

	// "github.com/Vesninovich/go-tasks/todos/task/inmemory"
	taskservice "github.com/Vesninovich/go-tasks/todos/task/service"
	tasksql "github.com/Vesninovich/go-tasks/todos/task/sql"
)

// TODO: get from config
const dbURL = "postgres://gotodos:gotodos@db:5432/gotodos"
const host = "0.0.0.0:3000"

func main() {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL DB at URL %s\n%s", dbURL, err)
	}
	defer db.Close()
	err = db.PingContext(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL DB at URL %s\n%s", dbURL, err)
	}
	log.Printf("Connected to PostgreSQL DB at URL %s\n", dbURL)

	// taskRepo := inmemory.New()
	taskRepo := tasksql.New(db)
	taskService := taskservice.New(taskRepo)
	taskServer := taskhttp.New(taskService)

	log.Printf("Starting server at host %s\n", host)
	_, err = httpserver.StartServer(host, "/api/v1", taskServer)
	if err != nil {
		log.Fatalf("Failed to start tasks server at host %s\n%s", host, err)
	}
	// log.Printf("Started tasks server at %s", tasksServer.Addr)
}
