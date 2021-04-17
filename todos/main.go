package main

import (
	"log"

	"github.com/Vesninovich/go-tasks/todos/http_server"
	task_http "github.com/Vesninovich/go-tasks/todos/task/http"
	"github.com/Vesninovich/go-tasks/todos/task/inmemory"
	task_service "github.com/Vesninovich/go-tasks/todos/task/service"
)

func main() {
	taskRepo := inmemory.New()
	taskService := task_service.New(taskRepo)
	taskServer := task_http.New(taskService)
	_, err := http_server.StartServer("localhost:3000", "/api/v1", taskServer)
	if err != nil {
		log.Fatalf("Failed to start tasks server: %s", err.Error())
	}
	// log.Printf("Started tasks server at %s", tasksServer.Addr)
}
