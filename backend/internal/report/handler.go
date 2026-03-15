package report

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rupi/timetracking/internal/auth"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	from, to := parseTimeRange(r)
	result, err := h.svc.Summary(r.Context(), userID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate summary")
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) Daily(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	from, to := parseTimeRange(r)
	tz := r.URL.Query().Get("timezone")
	if tz == "" {
		tz = "UTC"
	}
	result, err := h.svc.Daily(r.Context(), userID, from, to, tz)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate daily report")
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) Weekly(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	from, to := parseTimeRange(r)
	result, err := h.svc.Weekly(r.Context(), userID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate weekly report")
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func parseTimeRange(r *http.Request) (time.Time, time.Time) {
	q := r.URL.Query()
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now()

	if v := q.Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from = t
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			to = t
		}
	}
	return from, to
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
