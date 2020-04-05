package rest

import (
	"github.com/BackAged/domain/task"
	"net/http"
)

// AlbumHandler interface for the album handlers.
type Handler interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	GetUserTask(http.ResponseWriter, *http.Request
}


type hander struct {
	tskSvc Task.Service
}

// NewHandler will instantiate the handlers.
func NewHandler(tskSvc Task.Service) AlbumHandler {
	return &handler{tskSvc}
}

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	err := h.tskSvc.Create(r.Context())
	if err != nil {
		InvalidRequest(w, r, err, errFindAll, http.StatusUnprocessableEntity)
		return
	}

	ToJSON(w, http.StatusOK, &albums)
}