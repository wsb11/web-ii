package model

import (
	"errors"
	"html"
	"net/url"
	"strings"
	"time"
)

type Aluno struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Foto  string `json:"foto,omitempty"`
	Turma string `json:"turma,omitempty"`
}

type AlunoInput struct {
	Nome  string `json:"nome"`
	Foto  string `json:"foto,omitempty"`
	Turma string `json:"turma,omitempty"`
}

type AlunoComFotos struct {
	Aluno
	Fotos []Foto `json:"fotos"`
}

type Foto struct {
	ID      int    `json:"id"`
	AlunoID int    `json:"aluno_id"`
	URL     string `json:"url"`
	Legenda string `json:"legenda,omitempty"`
}

type FotoInput struct {
	URL     string `json:"url"`
	Legenda string `json:"legenda,omitempty"`
}

type Evento struct {
	ID        int    `json:"id"`
	Titulo    string `json:"titulo"`
	Descricao string `json:"descricao"`
	Data      string `json:"data"`
	ImagemURL string `json:"imagem_url,omitempty"`
}

type EventoInput struct {
	Titulo    string `json:"titulo"`
	Descricao string `json:"descricao"`
	Data      string `json:"data"`
	ImagemURL string `json:"imagem_url,omitempty"`
}

type Admin struct {
	ID        int
	Usuario   string
	SenhaHash string
	Role      string
}

type RefreshToken struct {
	ID        int
	AdminID   int
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func ValidarAluno(input AlunoInput) (AlunoInput, error) {
	input.Nome = sanitize(input.Nome)
	input.Foto = strings.TrimSpace(input.Foto)
	input.Turma = sanitize(input.Turma)

	if len(input.Nome) < 3 {
		return input, errors.New("nome deve ter pelo menos 3 caracteres")
	}
	if len(input.Nome) > 120 {
		return input, errors.New("nome deve ter no maximo 120 caracteres")
	}
	if len(input.Turma) > 30 {
		return input, errors.New("turma deve ter no maximo 30 caracteres")
	}
	if input.Foto != "" && !urlHTTPValida(input.Foto) {
		return input, errors.New("foto deve ser uma URL http ou https valida")
	}
	return input, nil
}

func ValidarFoto(input FotoInput) (FotoInput, error) {
	input.URL = strings.TrimSpace(input.URL)
	input.Legenda = sanitize(input.Legenda)

	if !urlHTTPValida(input.URL) {
		return input, errors.New("url da foto deve ser http ou https valida")
	}
	if len(input.Legenda) > 180 {
		return input, errors.New("legenda deve ter no maximo 180 caracteres")
	}
	return input, nil
}

func ValidarEvento(input EventoInput) (EventoInput, error) {
	input.Titulo = sanitize(input.Titulo)
	input.Descricao = sanitize(input.Descricao)
	input.Data = strings.TrimSpace(input.Data)
	input.ImagemURL = strings.TrimSpace(input.ImagemURL)

	if len(input.Titulo) < 3 {
		return input, errors.New("titulo deve ter pelo menos 3 caracteres")
	}
	if len(input.Titulo) > 160 {
		return input, errors.New("titulo deve ter no maximo 160 caracteres")
	}
	if len(input.Descricao) < 5 {
		return input, errors.New("descricao deve ter pelo menos 5 caracteres")
	}
	if _, err := time.Parse("2006-01-02", input.Data); err != nil {
		return input, errors.New("data deve estar no formato YYYY-MM-DD")
	}
	if input.ImagemURL != "" && !urlHTTPValida(input.ImagemURL) {
		return input, errors.New("imagem_url deve ser uma URL http ou https valida")
	}
	return input, nil
}

func sanitize(value string) string {
	return html.EscapeString(strings.TrimSpace(value))
}

func urlHTTPValida(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
