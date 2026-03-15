package timer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rupi/timetracking/internal/auth"
	"github.com/rupi/timetracking/internal/domain"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	entry, err := h.svc.GetActive(r.Context(), userID)
	if err != nil {
		respondJSON(w, http.StatusOK, nil)
		return
	}
	respondJSON(w, http.StatusOK, entry)
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	entry, err := h.svc.Start(r.Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondError(w, http.StatusBadRequest, "category_id is required")
		case domain.ErrActiveTimerExists:
			respondError(w, http.StatusConflict, "an active timer already exists")
		default:
			respondError(w, http.StatusInternalServerError, "failed to start timer")
		}
		return
	}
	respondJSON(w, http.StatusCreated, entry)
}

func (h *Handler) Pause(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	entry, err := h.svc.Pause(r.Context(), userID)
	if err != nil {
		handleTimerError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, entry)
}

func (h *Handler) Resume(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	entry, err := h.svc.Resume(r.Context(), userID)
	if err != nil {
		handleTimerError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, entry)
}

func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	entry, err := h.svc.Stop(r.Context(), userID)
	if err != nil {
		handleTimerError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, entry)
}

func (h *Handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid entry id")
		return
	}
	entry, err := h.svc.GetEntry(r.Context(), userID, entryID)
	if err != nil {
		respondError(w, http.StatusNotFound, "entry not found")
		return
	}
	respondJSON(w, http.StatusOK, entry)
}

func (h *Handler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid entry id")
		return
	}
	var body struct {
		Description string    `json:"description"`
		CategoryID  uuid.UUID `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	entry, err := h.svc.UpdateEntry(r.Context(), userID, entryID, body.Description, body.CategoryID)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "entry not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update entry")
		return
	}
	respondJSON(w, http.StatusOK, entry)
}

func (h *Handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid entry id")
		return
	}
	if err := h.svc.DeleteEntry(r.Context(), userID, entryID); err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "entry not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete entry")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListEntries(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())

	q := r.URL.Query()
	filter := ListFilter{
		Page:  1,
		Limit: 20,
	}

	if v := q.Get("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			filter.From = &t
		}
	}
	if v := q.Get("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			filter.To = &t
		}
	}
	if v := q.Get("category_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			filter.CategoryID = &id
		}
	}
	if v := q.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Page = n
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Limit = n
		}
	}

	entries, total, err := h.svc.ListEntries(r.Context(), userID, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list entries")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"entries": entries,
		"total":   total,
		"page":    filter.Page,
		"limit":   filter.Limit,
	})
}

func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())

	flusher, ok := w.(http.Flusher)
	if !ok {
		respondError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			elapsed, entry, err := h.svc.GetElapsed(r.Context(), userID)
			if err != nil {
				fmt.Fprintf(w, "event: idle\ndata: {}\n\n")
				flusher.Flush()
				continue
			}
			data, _ := json.Marshal(map[string]interface{}{
				"entry_id":    entry.ID,
				"state":       entry.State,
				"elapsed_sec": elapsed,
				"category_id": entry.CategoryID,
			})
			fmt.Fprintf(w, "event: tick\ndata: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func handleTimerError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrNotFound:
		respondError(w, http.StatusNotFound, "no active timer")
	case domain.ErrInvalidTransition:
		respondError(w, http.StatusConflict, "invalid timer state transition")
	default:
		respondError(w, http.StatusInternalServerError, "timer operation failed")
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
