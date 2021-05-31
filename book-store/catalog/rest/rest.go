package rest

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	bookservice "github.com/Vesninovich/go-tasks/book-store/catalog/book/service"
	_ "github.com/Vesninovich/go-tasks/book-store/catalog/docs" // generated docs
	"github.com/Vesninovich/go-tasks/book-store/common/book"
	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

var ctx = context.Background()

// Server of orders
type Server struct {
	service *bookservice.BookService
	baseURL string
	host    string
}

type apiModel struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Author     string   `json:"author"`
	Categories []string `json:"categories"`
}

type createAPIModel struct {
	Name   string `json:"name"`
	Author struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Categories []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		ParentID string `json:"parentID"`
	} `json:"categories"`
}

// New creates Server
func New(host, baseURL string, s *bookservice.BookService) *Server {
	return &Server{
		service: s,
		baseURL: baseURL,
		host:    host,
	}
}

// Start builds HTTP server for application and attempts to start it on given host.
// Created server serves requests starting from given `baseURL`.
func (s *Server) Start() error {
	serveMux := http.NewServeMux()
	s.handleTaskEndpoints(serveMux)
	var server http.Server
	server.Handler = serveMux
	server.Addr = s.host
	err := server.ListenAndServe()
	return err
}

func (s *Server) handleTaskEndpoints(serveMux *http.ServeMux) {
	serveMux.HandleFunc(s.baseURL, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.getBooks(w, r)
		case http.MethodPost:
			s.createBook(w, r)
		default:
			writeNotFound(w)
		}
	})

	// validPath := regexp.MustCompile(s.baseURL + "/" + uuid.REGEX + "$")
	// serveMux.HandleFunc(s.baseURL+"/", func(w http.ResponseWriter, r *http.Request) {
	// 	m := validPath.FindStringSubmatch(r.URL.Path)
	// 	if m == nil {
	// 		writeNotFound(w)
	// 		return
	// 	}
	// 	switch r.Method {
	// 	case http.MethodPut:
	// 		s.updateDescription(w, r)
	// 	case http.MethodDelete:
	// 		s.removeOrder(w, r)
	// 	default:
	// 		writeNotFound(w)
	// 	}
	// })

	// TODO: move to separate server
	serveMux.HandleFunc(s.baseURL+"/swagger/", httpSwagger.Handler(httpSwagger.URL(s.baseURL+"/swagger/doc.json")))
}

// getBooks godoc
// @Summary get books
// @Description get books according to query
// @Tags Book
// @Produce json
// @Param from query string false "results start"
// @Param count query string false "results count"
// @Param id query string false "book id"
// @Param author query string false "author id"
// @Param categories query []string false "category ids"
// @Success 200 {object} []apiModel "results"
// @Failure 400 {string} string "malformed query"
// @Failure 500 {string} string "internal error"
// @Router /book [get]
func (s *Server) getBooks(w http.ResponseWriter, r *http.Request) {
	from, count, query, err := parseQuery(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	books, err := s.service.GetBooks(ctx, from, count, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
	models := make([]apiModel, len(books))
	for i, b := range books {
		models[i] = toResponse(b)
	}
	res, err := json.Marshal(models)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

// createBook godoc
// @Summary create book
// @Description create book
// @Tags Book
// @Accept json
// @Produce json
// @Param order body createAPIModel true "book data"
// @Success 200 {object} book.Book "created book"
// @Failure 400 {string} string "malformed data"
// @Failure 404 {string} string "nested author or category not found"
// @Failure 500 {string} string "internal error"
// @Router /book [post]
func (s *Server) createBook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var data createAPIModel
	err = json.Unmarshal(body, &data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	aut := book.Author{Name: data.Author.Name}
	if data.Author.ID != "" {
		id, err := uuid.FromString(data.Author.ID)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		aut.ID = id
	}
	cats := make([]book.Category, len(data.Categories))
	for i, cat := range data.Categories {
		cats[i].Name = cat.Name
		if cat.ID != "" {
			id, err := uuid.FromString(cat.ID)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			cats[i].ID = id
		}
		if cat.ParentID != "" {
			id, err := uuid.FromString(cat.ParentID)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			cats[i].ParentID = id
		}
	}
	b, err := s.service.CreateBook(ctx, data.Name, aut, cats)
	writeResponse(w, b, err)
}

func getUUIDFromURL(path string) (uuid.UUID, error) {
	parts := strings.Split(path, "/")
	return uuid.FromString(parts[len(parts)-1])
}

// TODO: split (or not)
func parseQuery(params url.Values) (from, count uint, query book.Query, err error) {
	var val uint64
	var id uuid.UUID
	param := params.Get("from")
	if param != "" {
		val, err = strconv.ParseUint(param, 10, 64)
		if err != nil {
			return
		}
		from = uint(val)
	}
	param = params.Get("count")
	if param != "" {
		val, err = strconv.ParseUint(param, 10, 64)
		if err != nil {
			return
		}
		count = uint(val)
	}
	param = params.Get("id")
	if param != "" {
		id, err = uuid.FromString(param)
		if err != nil {
			return
		}
		query.ID = id
	}
	param = params.Get("author")
	if param != "" {
		id, err = uuid.FromString(param)
		if err != nil {
			return
		}
		query.Author = id
	}
	param = params.Get("categories")
	if param != "" {
		cats := strings.Split(param, ",")
		query.Categories = make([]uuid.UUID, len(cats))
		for i, cat := range cats {
			id, err = uuid.FromString(cat)
			if err != nil {
				return
			}
			query.Categories[i] = id
		}
	}
	return
}

func writeResponse(w http.ResponseWriter, b book.Book, err error) {
	if err != nil {
		switch err.(type) {
		case *commonerrors.NotFound:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	res, err := json.Marshal(toResponse(b))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

// I do not like the message written by http.NotFound() method
func writeNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func toResponse(b book.Book) apiModel {
	cats := make([]string, len(b.Categories))
	if len(b.Categories) != 0 {
		for i, cat := range b.Categories {
			cats[i] = cat.ID.String()
		}
	}
	return apiModel{
		ID:         b.ID.String(),
		Name:       b.Name,
		Author:     b.Author.ID.String(),
		Categories: cats,
	}
}
