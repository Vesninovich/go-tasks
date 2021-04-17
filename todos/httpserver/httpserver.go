package httpserver

import (
	"net/http"

	"github.com/Vesninovich/go-tasks/todos/task"
)

// StartServer builds HTTP server for application and attempts to start it on given host.
// Created server serves requests starting from given `baseURL`.
func StartServer(host, baseURL string, taskServer task.TasksServer) (*http.Server, error) {
	serveMux := http.NewServeMux()
	handleTaskEndpoints(serveMux, taskServer, baseURL+"/task")
	var server http.Server
	server.Handler = serveMux
	server.Addr = host
	err := server.ListenAndServe()
	return &server, err
}

func handleTaskEndpoints(serveMux *http.ServeMux, taskServer task.TasksServer, baseURL string) {
	serveMux.HandleFunc(baseURL, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskServer.GetTasks(w, r)
		case http.MethodPost:
			taskServer.PostTask(w, r)
		}
	})
}
