package rest

import (
	"encoding/json"
	"net/http"

	"github.com/BackAged/go-hexagonal-architecture/domain/task"
	"github.com/go-chi/chi"
)

// TaskRouter contains all routes for albums service.
func TaskRouter(h Handler) http.Handler {
	router := chi.NewRouter()

	router.Get("/{id}", h.Get)
	router.Post("/create", h.Create)
	router.Post("/{userid}", h.GetUserTask)

	return router
}

// Handler interface for the task handlers.
type Handler interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	GetUserTask(http.ResponseWriter, *http.Request)
}

type handler struct {
	tskSvc task.Service
}

// NewHandler will instantiate the handler
func NewHandler(tskSvc task.Service) Handler {
	return &handler{tskSvc: tskSvc}
}

// Create handler
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tsk := &task.Task{}
	err := json.NewDecoder(r.Body).Decode(&tsk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.tskSvc.Create(r.Context(), tsk); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(tsk.ID)
}

// Get handler
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	tskID := chi.URLParam(r, "id")
	tsk, err := h.tskSvc.Get(r.Context(), tskID)
	if err != nil {
		serializeJson(w, 500, "something went wrong")
	}
	if tsk == nil {
		serializeJson(w, 400, "Not Found")
	}

	serializeJson(w, 200, tsk)
}

// GetUserTask handler
func (h *handler) GetUserTask(w http.ResponseWriter, r *http.Request) {

}
