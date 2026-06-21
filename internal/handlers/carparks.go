package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bryanngzh/parklah-go/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type batchRequest struct {
	Codes []string `json:"codes"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func queryFloat(r *http.Request, key string) (float64, bool) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	return v, err == nil
}

func queryInt(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func validSource(source string) bool {
	return source == "ura" || source == "hdb"
}

func requireSource(w http.ResponseWriter, r *http.Request) (string, bool) {
	source := r.URL.Query().Get("source")
	if source == "" {
		writeError(w, http.StatusBadRequest, "source query param is required (ura or hdb)")
		return "", false
	}
	if !validSource(source) {
		writeError(w, http.StatusBadRequest, "source must be ura or hdb")
		return "", false
	}
	return source, true
}

func GetNearby(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat, ok1 := queryFloat(r, "lat")
		lon, ok2 := queryFloat(r, "lon")
		if !ok1 || !ok2 {
			writeError(w, http.StatusBadRequest, "lat and lon are required")
			return
		}

		radius, ok := queryFloat(r, "radius")
		if !ok {
			radius = 600
		}
		if radius > 2000 {
			radius = 2000
		}

		limit := queryInt(r, "limit", 20)
		if limit > 50 {
			limit = 50
		}

		results, meta, err := services.GetNearby(r.Context(), pool, lat, lon, radius, limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch nearby carparks")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"data": results, "meta": meta})
	}
}

func GetBatch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat, ok1 := queryFloat(r, "lat")
		lon, ok2 := queryFloat(r, "lon")
		if !ok1 || !ok2 {
			writeError(w, http.StatusBadRequest, "lat and lon are required")
			return
		}

		var body batchRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Codes) == 0 {
			writeError(w, http.StatusBadRequest, "body must be {\"codes\": [...]}")
			return
		}

		results, err := services.GetBatch(r.Context(), pool, lat, lon, body.Codes)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch batch carparks")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"data": results, "meta": map[string]any{"count": len(results)}})
	}
}

func GetCarpark(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		source, ok := requireSource(w, r)
		if !ok {
			return
		}

		detail, err := services.GetCarparkDetail(r.Context(), pool, code, source)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch carpark")
			return
		}
		if detail == nil {
			writeError(w, http.StatusNotFound, "carpark not found")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"data": detail})
	}
}

func GetAvailability(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		source, ok := requireSource(w, r)
		if !ok {
			return
		}

		avail, err := services.GetAvailability(r.Context(), pool, code, source)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch availability")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"data": avail})
	}
}

func GetRates(pool *pgxpool.Pool, phDates map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		source, ok := requireSource(w, r)
		if !ok {
			return
		}

		rates, err := services.GetRates(r.Context(), pool, code, source, phDates)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch rates")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"data": rates})
	}
}
