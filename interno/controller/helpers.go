package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const maxJSONBody = 1 << 20

type errorResponse struct {
	Erro string `json:"erro"`
}

func responderJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func responderErro(w http.ResponseWriter, status int, message string) {
	responderJSON(w, status, errorResponse{Erro: message})
}

func decodificarJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			responderErro(w, http.StatusBadRequest, "corpo JSON obrigatorio")
			return false
		}
		responderErro(w, http.StatusBadRequest, "JSON invalido")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		responderErro(w, http.StatusBadRequest, "JSON deve conter apenas um objeto")
		return false
	}
	return true
}

func idDaURL(w http.ResponseWriter, r *http.Request, nome string) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, nome))
	if err != nil || id <= 0 {
		responderErro(w, http.StatusBadRequest, "ID invalido")
		return 0, false
	}
	return id, true
}
