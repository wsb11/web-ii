package controller

import (
	"net/http"
	"time"

	"anuario/interno/auth"
	"anuario/interno/model"
)

type loginRequest struct {
	Usuario string `json:"usuario"`
	Senha   string `json:"senha"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input loginRequest
	if !decodificarJSON(w, r, &input) {
		return
	}
	if input.Usuario == "" || input.Senha == "" {
		responderErro(w, http.StatusBadRequest, "usuario e senha sao obrigatorios")
		return
	}

	admin, encontrado, err := store.BuscarAdminPorUsuario(r.Context(), input.Usuario)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao autenticar")
		return
	}
	if !encontrado || !auth.VerifyPassword(input.Senha, admin.SenhaHash) {
		responderErro(w, http.StatusUnauthorized, "credenciais invalidas")
		return
	}

	responderTokens(w, r, admin)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var input refreshRequest
	if !decodificarJSON(w, r, &input) {
		return
	}
	if input.RefreshToken == "" {
		responderErro(w, http.StatusBadRequest, "refresh_token obrigatorio")
		return
	}

	hash := auth.HashRefreshToken(input.RefreshToken)
	token, encontrado, err := store.BuscarRefreshToken(r.Context(), hash)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao validar refresh token")
		return
	}
	if !encontrado {
		responderErro(w, http.StatusUnauthorized, "refresh token invalido")
		return
	}

	admin, encontrado, err := store.BuscarAdminPorID(r.Context(), token.AdminID)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao buscar administrador")
		return
	}
	if !encontrado {
		responderErro(w, http.StatusUnauthorized, "administrador nao encontrado")
		return
	}

	if err := store.RevogarRefreshToken(r.Context(), hash); err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao rotacionar refresh token")
		return
	}
	responderTokens(w, r, admin)
}

func responderTokens(w http.ResponseWriter, r *http.Request, admin model.Admin) {
	accessToken, err := tokenMaker.Generate(admin)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao gerar JWT")
		return
	}
	refreshToken, refreshHash, err := auth.NewRefreshToken()
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao gerar refresh token")
		return
	}
	if err := store.CriarRefreshToken(r.Context(), admin.ID, refreshHash, time.Now().Add(refreshTTL)); err != nil {
		responderErro(w, http.StatusInternalServerError, "erro ao salvar refresh token")
		return
	}

	responderJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(tokenMaker.AccessTTL().Seconds()),
	})
}
