package sites

import (
	"log"
	"net/http"
	"uptime-monitor/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.service.ListSites(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 2. Return JSON in an HTTP response
	json.Write(w, http.StatusOK, *sites)
}
