package task_http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Vesninovich/go-tasks/todos/common"
	"github.com/Vesninovich/go-tasks/todos/task"
	task_service "github.com/Vesninovich/go-tasks/todos/task/service"
)

type HttpServer struct {
	service *task_service.Service
}

func New(service *task_service.Service) *HttpServer {
	return &HttpServer{service}
}

func (s *HttpServer) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.service.GetAll(context.Background())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	res, err := json.Marshal(prepareTasks(tasks))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

func (s *HttpServer) PostTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var data createTaskApiModel
	err = json.Unmarshal(body, &data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	task, err := s.service.CreateTask(context.Background(), data.Name, data.Description, int64(data.DueDate))
	if err != nil {
		switch err.(type) {
		case *common.InvalidInputError:
			writeError(w, http.StatusBadRequest, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprint(task.ID)))
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

type taskApiModel struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DueDate     int64  `json:"dueDate,omitempty"`
	Status      string `json:"status,omitempty"`
}

type createTaskApiModel struct {
	Name        string
	Description string
	DueDate     int
}

func prepareTasks(tasks []task.Task) []taskApiModel {
	t := make([]taskApiModel, len(tasks))
	for i, task := range tasks {
		t[i] = taskToApiModel(task)
	}
	return t
}

func taskToApiModel(t task.Task) taskApiModel {
	return taskApiModel{
		t.ID,
		t.Name,
		t.Description,
		t.DueDate.Unix(),
		getStatusText(t.Status),
	}
}

func getStatusText(s task.Status) string {
	switch s {
	case task.New:
		return "new"
	case task.InProgress:
		return "in-progress"
	case task.Cancelled:
		return "cancelled"
	case task.Done:
		return "done"
	case task.Overdue:
		return "overdue"
	default:
		return ""
	}
}
