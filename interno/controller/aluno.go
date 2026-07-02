package controller

import (
	"net/http"

	"anuario/interno/model"
)

func ListarAlunos(w http.ResponseWriter, r *http.Request) {
	alunos, err := store.ListarAlunos(r.Context())
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao listar alunos")
		return
	}
	responderJSON(w, http.StatusOK, alunos)
}

func ObterAluno(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}

	aluno, encontrado, err := store.BuscarAlunoComFotos(r.Context(), id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao buscar aluno")
		return
	}
	if !encontrado {
		responderErro(w, http.StatusNotFound, "aluno nao encontrado")
		return
	}
	responderJSON(w, http.StatusOK, aluno)
}

func CriarAluno(w http.ResponseWriter, r *http.Request) {
	var input model.AlunoInput
	if !decodificarJSON(w, r, &input) {
		return
	}

	input, err := model.ValidarAluno(input)
	if err != nil {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}

	alunoCriado, err := store.AdicionarAluno(r.Context(), input)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao criar aluno")
		return
	}
	responderJSON(w, http.StatusCreated, alunoCriado)
}

func AtualizarAluno(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}

	var input model.AlunoInput
	if !decodificarJSON(w, r, &input) {
		return
	}
	input, err := model.ValidarAluno(input)
	if err != nil {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}

	aluno, encontrado, err := store.AtualizarAluno(r.Context(), id, input)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao atualizar aluno")
		return
	}
	if !encontrado {
		responderErro(w, http.StatusNotFound, "aluno nao encontrado")
		return
	}
	responderJSON(w, http.StatusOK, aluno)
}

func RemoverAluno(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}

	removido, err := store.RemoverAluno(r.Context(), id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao remover aluno")
		return
	}
	if !removido {
		responderErro(w, http.StatusNotFound, "aluno nao encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ListarFotosDoAluno(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}
	if _, encontrado, err := store.BuscarAlunoPorID(r.Context(), id); err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao buscar aluno")
		return
	} else if !encontrado {
		responderErro(w, http.StatusNotFound, "aluno nao encontrado")
		return
	}

	fotos, err := store.ListarFotosPorAluno(r.Context(), id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao listar fotos")
		return
	}
	responderJSON(w, http.StatusOK, fotos)
}

func CriarFotoDoAluno(w http.ResponseWriter, r *http.Request) {
	id, ok := idDaURL(w, r, "id")
	if !ok {
		return
	}

	var input model.FotoInput
	if !decodificarJSON(w, r, &input) {
		return
	}
	input, err := model.ValidarFoto(input)
	if err != nil {
		responderErro(w, http.StatusBadRequest, err.Error())
		return
	}

	foto, encontrado, err := store.CriarFoto(r.Context(), id, input)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao criar foto")
		return
	}
	if !encontrado {
		responderErro(w, http.StatusNotFound, "aluno nao encontrado")
		return
	}
	responderJSON(w, http.StatusCreated, foto)
}
