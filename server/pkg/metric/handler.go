package metric

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct {
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, URL, h.heartbeat)
}

// Heartbeat
// @Summary Heartbeat metric
// @Tags Metric
// @Success 204
// @Failure 400
// @Router /api/heartbeat [get]
func (h *Handler) heartbeat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
}
