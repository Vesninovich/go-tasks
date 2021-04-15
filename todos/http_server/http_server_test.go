package http_server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTasksHttp(t *testing.T) {
	res := getAll(t)
	if res != "[]" {
		t.Errorf("Expected to get empty JSON array '[]', got\n\t%s", res)
	}
}

func getAll(t *testing.T) string {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	readAll(rec, req)

	res := rec.Result()
	body, err := ioutil.ReadAll(res.Body)

	status := res.StatusCode
	contentType := strings.ToLower(res.Header.Get("Content-Type"))
	switch {
	case err != nil:
		t.Errorf("Got error reading response body:\n\t%s", err.Error())
	case status != http.StatusOK:
		t.Errorf("Expected to get status %d, got %d", http.StatusOK, status)
	case contentType != "application/json":
		t.Errorf("Expected to get response of content type 'application/json', got %s", contentType)
	}
	return string(body)
}
