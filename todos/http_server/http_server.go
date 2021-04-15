package http_server

import (
	"io"
	"net/http"
)

func readAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "[]")
}

func StartServer(host, baseUrl string) (*http.Server, error) {
	// http.HandleFunc("/", readAll)
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(baseUrl, readAll)
	// server := &http.Server{}
	var server http.Server
	server.Handler = serveMux
	server.Addr = host
	err := server.ListenAndServe()
	return &server, err
}
