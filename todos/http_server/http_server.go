package http_server

import (
	"net/http"

	"github.com/Vesninovich/go-tasks/todos/task"
)

func StartServer(host, baseUrl string, taskServer task.TasksServer) (*http.Server, error) {
	serveMux := http.NewServeMux()
	handleTaskEndpoints(serveMux, taskServer, baseUrl+"/task")
	var server http.Server
	server.Handler = serveMux
	server.Addr = host
	err := server.ListenAndServe()
	return &server, err
}

func handleTaskEndpoints(serveMux *http.ServeMux, taskServer task.TasksServer, baseUrl string) {
	serveMux.HandleFunc(baseUrl, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskServer.GetTasks(w, r)
		case http.MethodPost:
			taskServer.PostTask(w, r)
		}
	})
}
