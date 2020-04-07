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
	router.Get("/user/{userid}", h.GetUserTask)

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

type createDTO struct {
	UserID      string          `json:"userID"`
	Topic       string          `json:"topic"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	SubTasks    []*task.SubTask `json:"sub_task"`
}

// Create handler
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tskDTO := &createDTO{}
	if err := json.NewDecoder(r.Body).Decode(&tskDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tsk := &task.Task{
		UserID:      tskDTO.UserID,
		Topic:       tskDTO.Topic,
		Description: tskDTO.Description,
		Status:      task.Status(tskDTO.Status),
		SubTasks:    tskDTO.SubTasks,
	}

	if err := h.tskSvc.Create(r.Context(), tsk); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(tsk.ID)
}

// Get handler
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tskID := chi.URLParam(r, "id")
	tsk, err := h.tskSvc.Get(r.Context(), tskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tsk == nil {
		http.Error(w, "Notfound", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tsk)

}

// GetUserTask handler
func (h *handler) GetUserTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usrID := chi.URLParam(r, "userid")
	tsk, err := h.tskSvc.GetUserTask(r.Context(), usrID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tsk == nil {
		http.Error(w, "Notfound", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tsk)
}
