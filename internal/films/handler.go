package films

import (
	"net/http"

	"example.com/filmserver/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{service: NewService()}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	films, err := h.service.List()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, films)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	film, err := h.service.Get(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	response.JSON(w, http.StatusOK, film)
}
