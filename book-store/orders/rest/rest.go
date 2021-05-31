package rest

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
	_ "github.com/Vesninovich/go-tasks/book-store/orders/docs" // generated docs
	"github.com/Vesninovich/go-tasks/book-store/orders/order"
	orderservice "github.com/Vesninovich/go-tasks/book-store/orders/order/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

var ctx = context.Background()

// Server of orders
type Server struct {
	service *orderservice.Service
	baseURL string
	host    string
}

type apiModel struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	BookID      string `json:"bookID"`
}

type createAPIModel struct {
	Description string `json:"description"`
	BookID      string `json:"bookID"`
}

type descUpdAPIModel struct {
	Description string `json:"description"`
}

// New creates Server
func New(host, baseURL string, s *orderservice.Service) *Server {
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
		case http.MethodPost:
			s.createOrder(w, r)
		default:
			writeNotFound(w)
		}
	})

	validPath := regexp.MustCompile(s.baseURL + "/" + uuid.REGEX + "$")
	serveMux.HandleFunc(s.baseURL+"/", func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			writeNotFound(w)
			return
		}
		switch r.Method {
		case http.MethodGet:
			s.getOrder(w, r)
		case http.MethodPut:
			s.updateDescription(w, r)
		case http.MethodDelete:
			s.removeOrder(w, r)
		default:
			writeNotFound(w)
		}
	})

	// TODO: move to separate server
	serveMux.HandleFunc(s.baseURL+"/swagger/", httpSwagger.Handler(httpSwagger.URL(s.baseURL+"/swagger/doc.json")))
}

// GetOrder godoc
// @Summary get order
// @Description get order by id
// @Tags Order
// @Produce json
// @Param id path string true "order id"
// @Success 200 {object} apiModel "requested order"
// @Failure 400 {string} string "malformed id"
// @Failure 404 {string} string "requested order not found"
// @Failure 500 {string} string "internal error"
// @Router /order/{id} [get]
func (s *Server) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getUUIDFromURL(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	o, err := s.service.GetOrder(ctx, id)
	writeResponse(w, o, err)
}

// CreateOrder godoc
// @Summary place order
// @Description place new book order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "order id"
// @Param order body createAPIModel true "order data"
// @Success 200 {object} apiModel "created order"
// @Failure 400 {string} string "malformed book id"
// @Failure 404 {string} string "requested book not found"
// @Failure 500 {string} string "internal error"
// @Router /order [post]
func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
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
	bID, err := uuid.FromString(data.BookID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	o, err := s.service.CreateOrder(ctx, order.CreateDTO{
		Description: data.Description,
		BookID:      bID,
	})
	writeResponse(w, o, err)
}

// UpdateDescription godoc
// @Summary update description
// @Description change order description
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "order id"
// @Param description body descUpdAPIModel true "new description"
// @Success 200 {object} apiModel "updated order"
// @Failure 400 {string} string "malformed order id or bad data"
// @Failure 404 {string} string "requested order not found"
// @Failure 500 {string} string "internal error"
// @Router /order/{id} [put]
func (s *Server) updateDescription(w http.ResponseWriter, r *http.Request) {
	id, err := getUUIDFromURL(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var data descUpdAPIModel
	err = json.Unmarshal(body, &data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	o, err := s.service.UpdateDescription(ctx, order.Order{
		ID:          id,
		Description: data.Description,
	})
	writeResponse(w, o, err)
}

// RemoveOrder godoc
// @Summary remove order
// @Description remove order
// @Tags Order
// @Produce json
// @Param order path string true "order id"
// @Success 200 {object} apiModel "removed order"
// @Failure 400 {string} string "malformed order id"
// @Failure 404 {string} string "requested order not found"
// @Failure 500 {string} string "internal error"
// @Router /order/{id} [delete]
func (s *Server) removeOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getUUIDFromURL(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	o, err := s.service.RemoveOrder(ctx, id)
	writeResponse(w, o, err)
}

func getUUIDFromURL(url string) (uuid.UUID, error) {
	parts := strings.Split(url, "/")
	return uuid.FromString(parts[len(parts)-1])
}

func writeResponse(w http.ResponseWriter, o order.Order, err error) {
	if err != nil {
		switch err.(type) {
		case *commonerrors.NotFound:
			writeError(w, http.StatusNotFound, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}
	res, err := json.Marshal(orderToResponse(o))
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

func orderToResponse(o order.Order) apiModel {
	return apiModel{
		ID:          o.ID.String(),
		Description: o.Description,
		BookID:      o.Book.ID.String(),
	}
}
