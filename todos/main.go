package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/Vesninovich/go-tasks/todos/httpserver"
	taskhttp "github.com/Vesninovich/go-tasks/todos/task/http"

	// "github.com/Vesninovich/go-tasks/todos/task/inmemory"
	taskservice "github.com/Vesninovich/go-tasks/todos/task/service"
	tasksql "github.com/Vesninovich/go-tasks/todos/task/sql"
)

func main() {
	dbURL := buildDbURL()
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

	host := buildHost()
	log.Printf("Starting server at host %s\n", host)
	_, err = httpserver.StartServer(host, "/api/v1", taskServer)
	if err != nil {
		log.Fatalf("Failed to start tasks server at host %s\n%s", host, err)
	}
	// log.Printf("Started tasks server at %s", tasksServer.Addr)
}

func buildDbURL() (dbURL string) {
	db := os.Getenv("TODO_DB")
	host := os.Getenv("TODO_DB_HOST")
	port := os.Getenv("TODO_DB_PORT")
	user := os.Getenv("TODO_DB_USER")
	pwd, pwdSet := os.LookupEnv("TODO_DB_PWD")
	if db == "" {
		db = "gotodos"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "gotodos"
	}
	if !pwdSet {
		pwd = ":gotodos"
	} else if pwd != "" {
		pwd = ":" + pwd
	}
	return fmt.Sprintf("postgresql://%s%s@%s:%s/%s", user, pwd, host, port, db)
}

func buildHost() (host string) {
	host = os.Getenv("TODO_HOST")
	if host == "" {
		return "0.0.0.0:3000"
	}
	return
}
