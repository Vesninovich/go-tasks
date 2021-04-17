package taskhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Vesninovich/go-tasks/todos/common"
	"github.com/Vesninovich/go-tasks/todos/task"
	task_service "github.com/Vesninovich/go-tasks/todos/task/service"
)

// HTTPServer serves requests for Tasks
type HTTPServer struct {
	service *task_service.Service
}

// New creates new instance of HTTPServer
func New(service *task_service.Service) *HTTPServer {
	return &HTTPServer{service}
}

// GetTasks serves requests to read slice of tasks
func (s *HTTPServer) GetTasks(w http.ResponseWriter, r *http.Request) {
	from, count, err := parsePaginationQuery(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	tasks, err := s.service.Get(context.Background(), uint(from), uint(count))
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

// GetOneTask serves requests to get task by id
func (s *HTTPServer) GetOneTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	tsk, err := s.service.GetOne(context.Background(), id)
	if err != nil {
		switch err.(type) {
		case *common.NotFoundError:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	res, err := json.Marshal(taskToAPIModel(tsk))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

// PostTask serves requests to create new task
func (s *HTTPServer) PostTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var data createTaskAPIModel
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

// PutTask serves requests to update existing task
func (s *HTTPServer) PutTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var data updateTaskAPIModel
	err = json.Unmarshal(body, &data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	id, err := getIDFromURL(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	_, err = s.service.UpdateTask(context.Background(), id, data.Name, data.Description, int64(data.DueDate), data.Status)
	if err != nil {
		switch err.(type) {
		case *common.InvalidInputError:
			writeError(w, http.StatusBadRequest, err)
		case *common.NotFoundError:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteTask serves requests to get task by id
func (s *HTTPServer) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	err = s.service.Delete(context.Background(), id)
	if err != nil {
		switch err.(type) {
		case *common.NotFoundError:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func getIDFromURL(url string) (uint64, error) {
	parts := strings.Split(url, "/")
	return strconv.ParseUint(parts[len(parts)-1], 10, 64)
}

func parsePaginationQuery(query url.Values) (from uint64, count uint64, err error) {
	f := query.Get("from")
	c := query.Get("count")
	if f == "" {
		from = 0
	} else {
		from, err = strconv.ParseUint(f, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	if c == "" {
		count = 0
	} else {
		count, err = strconv.ParseUint(c, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	return
}

type taskAPIModel struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DueDate     int64  `json:"dueDate,omitempty"`
	Status      string `json:"status,omitempty"`
}

type createTaskAPIModel struct {
	Name        string
	Description string
	DueDate     int
}

type updateTaskAPIModel struct {
	createTaskAPIModel
	ID     uint64
	Status string
}

func prepareTasks(tasks []task.Task) []taskAPIModel {
	t := make([]taskAPIModel, len(tasks))
	for i, task := range tasks {
		t[i] = taskToAPIModel(task)
	}
	return t
}

func taskToAPIModel(t task.Task) taskAPIModel {
	return taskAPIModel{
		t.ID,
		t.Name,
		t.Description,
		t.DueDate.Unix(),
		t.Status.String(),
	}
}
