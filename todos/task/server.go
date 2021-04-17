package task

import "net/http"

// TasksServer interface represents objects that serve HTTP requests for Tasks
type TasksServer interface {
	GetTasks(w http.ResponseWriter, r *http.Request)
	PostTask(w http.ResponseWriter, r *http.Request)
}
