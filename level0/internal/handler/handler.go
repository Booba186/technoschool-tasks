package handler

import (
	"encoding/json"
	"github.com/Booba186/level0/internal/cache"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	cache *cache.Cache
}

func NewHandler(c *cache.Cache) *Handler {
	return &Handler{cache: c}
}

func (h *Handler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "UID заказа не указан", http.StatusBadRequest)
		return
	}

	order, found := h.cache.Get(uid)
	if !found {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
