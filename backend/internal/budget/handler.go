package budget

import (
	"encoding/json"
	"net/http"

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	rules, err := h.svc.List(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list budget rules")
		return
	}
	respondJSON(w, http.StatusOK, rules)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	rule, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		if err == domain.ErrInvalidInput {
			respondError(w, http.StatusBadRequest, "category_id and budget_sec > 0 are required")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create budget rule")
		return
	}
	respondJSON(w, http.StatusCreated, rule)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	ruleID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid rule id")
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	rule, err := h.svc.Update(r.Context(), userID, ruleID, req)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "budget rule not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update budget rule")
		return
	}
	respondJSON(w, http.StatusOK, rule)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	ruleID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid rule id")
		return
	}
	if err := h.svc.Delete(r.Context(), userID, ruleID); err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "budget rule not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete budget rule")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	tz := r.URL.Query().Get("timezone")
	if tz == "" {
		tz = "UTC"
	}
	statuses, err := h.svc.GetStatus(r.Context(), userID, tz)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get budget status")
		return
	}
	respondJSON(w, http.StatusOK, statuses)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
