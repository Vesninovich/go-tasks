package main

import (
	"log"

	"github.com/Vesninovich/go-tasks/todos/httpserver"
	taskhttp "github.com/Vesninovich/go-tasks/todos/task/http"
	"github.com/Vesninovich/go-tasks/todos/task/inmemory"
	taskservice "github.com/Vesninovich/go-tasks/todos/task/service"
)

func main() {
	taskRepo := inmemory.New()
	taskService := taskservice.New(taskRepo)
	taskServer := taskhttp.New(taskService)
	_, err := httpserver.StartServer("localhost:3000", "/api/v1", taskServer)
	if err != nil {
		log.Fatalf("Failed to start tasks server: %s", err.Error())
	}
	// log.Printf("Started tasks server at %s", tasksServer.Addr)
}
