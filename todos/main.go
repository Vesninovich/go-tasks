package main

import (
	"log"

	"github.com/Vesninovich/go-tasks/todos/http_server"
)

func main() {
	_, err := http_server.StartServer("localhost:3000", "/api/v1/task")
	if err != nil {
		log.Fatalf("Failed to start tasks server: %s", err.Error())
	}
	// log.Printf("Started tasks server at %s", tasksServer.Addr)
}
