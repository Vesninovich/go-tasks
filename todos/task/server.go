package task

import "net/http"

type TasksServer interface {
	GetTasks(w http.ResponseWriter, r *http.Request)
	PostTask(w http.ResponseWriter, r *http.Request)
}
