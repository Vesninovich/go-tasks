package taskhttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Vesninovich/go-tasks/todos/task/inmemory"
	task_service "github.com/Vesninovich/go-tasks/todos/task/service"
)

func TestPostValid(t *testing.T) {
	taskJSONs := []string{
		`{"name":"testA","dueDate":12345678}`,
		`{"name":"testB","description":"asd"}`,
		`{"name":"testC","description":"dsa","dueDate":87654321}`,
	}
	savedTasks := []string{
		`{"id":0,"name":"testA","dueDate":12345678,"status":"new"}`,
		`{"id":1,"name":"testB","description":"asd","status":"new"}`,
		`{"id":2,"name":"testC","description":"dsa","dueDate":87654321,"status":"new"}`,
	}
	tasksJSON := "[" + strings.Join(savedTasks, ",") + "]"

	s := createServer()
	for _, task := range taskJSONs {
		status, contentType, body := postTask(t, s, task)
		_, err := strconv.ParseUint(string(body), 10, 64)
		if err != nil {
			t.Errorf("Got error parsing uint from body: %s", err)
		}
		checkStatus(t, http.StatusCreated, status)
		checkContentType(t, "text/plain", contentType)
	}

	status, contentType, body := getTasks(t, s)
	checkStatus(t, http.StatusOK, status)
	checkContentType(t, "application/json", contentType)
	if body != tasksJSON {
		t.Errorf("Expected tasks JSON to be \n\t%s\ngot\n\t%s", tasksJSON, body)
	}
}

func TestPostInvalid(t *testing.T) {
	taskJSONs := []string{
		`{"name":"","dueDate":12345678}`,
		`{"description":"asd"}`,
		`{"name":"testC","description":"dsa","dueDate":-1234}`,
		`{"name":"","description":"dsa","dueDate":-1234}`,
	}

	s := createServer()
	for _, task := range taskJSONs {
		status, contentType, _ := postTask(t, s, task)
		checkStatus(t, http.StatusBadRequest, status)
		checkContentType(t, "text/plain", contentType)
	}
}

func getTasks(t *testing.T, s *HTTPServer) (status int, contentType, body string) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	s.GetTasks(rec, req)
	bodyRaw, err := ioutil.ReadAll(rec.Result().Body)
	if err != nil {
		t.Errorf("Got error reading response body: %s", err)
	}
	return rec.Result().StatusCode, rec.Result().Header.Get("Content-Type"), string(bodyRaw)
}

func postTask(t *testing.T, s *HTTPServer, task string) (status int, contentType, body string) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(task))
	rec := httptest.NewRecorder()

	s.PostTask(rec, req)
	bodyRaw, err := ioutil.ReadAll(rec.Result().Body)
	if err != nil {
		t.Errorf("Got error reading response body: %s", err)
	}
	return rec.Result().StatusCode, rec.Result().Header.Get("Content-Type"), string(bodyRaw)
}

func checkContentType(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected to get content type %s, got %s", expected, actual)
	}
}

func checkStatus(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected to get status %d, got %d", expected, actual)
	}
}

func createServer() *HTTPServer {
	return New(task_service.New(inmemory.New()))
}
