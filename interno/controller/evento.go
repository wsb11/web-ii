package controller

import (
	"net/http"

	"anuario/interno/model"
)

func ListarEventos(w http.ResponseWriter, r *http.Request) {
	eventos, err := store.ListarEventos(r.Context())
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao listar eventos")
		return
	}
	responderJSON(w, http.StatusOK, eventos)
}

func ObterEvento(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}
	evento, encontrado, err := store.BuscarEventoPorID(r.Context(), id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao buscar evento")
		return
	}
	if !encontrado {
		responderErro(w, http.StatusNotFound, "evento nao encontrado")
		return
	}
	responderJSON(w, http.StatusOK, evento)
}

func CriarEvento(w http.ResponseWriter, r *http.Request) {
	var input model.EventoInput
	if !decodificarJSON(w, r, &input) {
		return
	}
	input, err := model.ValidarEvento(input)
	if err != nil {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}

	evento, err := store.CriarEvento(r.Context(), input)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao criar evento")
		return
	}
	responderJSON(w, http.StatusCreated, evento)
}

func AtualizarEvento(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}

	var input model.EventoInput
	if !decodificarJSON(w, r, &input) {
		return
	}
	input, err := model.ValidarEvento(input)
	if err != nil {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}

	evento, encontrado, err := store.AtualizarEvento(r.Context(), id, input)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao atualizar evento")
		return
	}
	if !encontrado {
		responderErro(w, http.StatusNotFound, "evento nao encontrado")
		return
	}
	responderJSON(w, http.StatusOK, evento)
}

func RemoverEvento(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}
	removido, err := store.RemoverEvento(r.Context(), id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao remover evento")
		return
	}
	if !removido {
		responderErro(w, http.StatusNotFound, "evento nao encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
