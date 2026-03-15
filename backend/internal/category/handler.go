package category

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
	cats, err := h.svc.List(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list categories")
		return
	}
	respondJSON(w, http.StatusOK, cats)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cat, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		if err == domain.ErrInvalidInput {
			respondError(w, http.StatusBadRequest, "name is required")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create category")
		return
	}
	respondJSON(w, http.StatusCreated, cat)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	catID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	cat, err := h.svc.Get(r.Context(), userID, catID)
	if err != nil {
		respondError(w, http.StatusNotFound, "category not found")
		return
	}
	respondJSON(w, http.StatusOK, cat)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	catID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	cat, err := h.svc.Update(r.Context(), userID, catID, req)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "category not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update category")
		return
	}
	respondJSON(w, http.StatusOK, cat)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserID(r.Context())
	catID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	if err := h.svc.Delete(r.Context(), userID, catID); err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "category not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete category")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
