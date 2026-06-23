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
	json.Write(w, http.StatusOK, sites)
}

func (h *handler) EnqueuePollSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.service.EnqueuePollSites(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 2. Return JSON in an HTTP response
	json.Write(w, http.StatusOK, sites)
}

func (h *handler) PollSite(w http.ResponseWriter, r *http.Request) {

	var tempParams pollParams
	if err := json.Read(r, &tempParams); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.PollSite(r.Context(), tempParams)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Return JSON in an HTTP response
	json.Write(w, http.StatusOK, resp)
}

func (h *handler) AddSite(w http.ResponseWriter, r *http.Request) {
	var tempParams createAddParams
	if err := json.Read(r, &tempParams); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sites, err := h.service.AddSite(r.Context(), tempParams)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 2. Return JSON in an HTTP response
	json.Write(w, http.StatusOK, sites)
}

func (h *handler) RemoveSite(w http.ResponseWriter, r *http.Request) {
	var tempParams createIdParams
	if err := json.Read(r, &tempParams); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name, err := h.service.RemoveSite(r.Context(), tempParams)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 2. Return JSON in an HTTP response
	json.Write(w, http.StatusOK, name)
}

func (h *handler) FindSiteByID(w http.ResponseWriter, r *http.Request) {
	var tempParams createIdParams
	if err := json.Read(r, &tempParams); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	site, err := h.service.FindSitesByID(r.Context(), tempParams)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, site)
}
