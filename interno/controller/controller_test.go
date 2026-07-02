package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"anuario/interno/auth"
	"anuario/interno/model"
	"anuario/interno/repository"
)

func setupControllerTest(t *testing.T) {
	t.Helper()
	Configurar(repository.NewMemoryStore(), auth.MustJWTManager("test-secret-anuario-2026", time.Minute))
	ConfigurarRefreshTTL(24 * time.Hour)
}

func withID(req *http.Request, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", value)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	return req.WithContext(ctx)
}

func TestListarAlunos(t *testing.T) {
	setupControllerTest(t)
	req := httptest.NewRequest("GET", "/api/v1/alunos", nil)
	rec := httptest.NewRecorder()

	ListarAlunos(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatal("Content-Type deve ser application/json")
	}
}

func TestObterAlunoExistenteRetornaFotosAninhadas(t *testing.T) {
	setupControllerTest(t)
	req := withID(httptest.NewRequest("GET", "/api/v1/alunos/1", nil), "1")
	rec := httptest.NewRecorder()

	ObterAluno(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}
	var aluno model.AlunoComFotos
	if err := json.NewDecoder(rec.Body).Decode(&aluno); err != nil {
		t.Fatal("resposta nao e JSON valido")
	}
	if aluno.ID != 1 || len(aluno.Fotos) == 0 {
		t.Fatalf("esperava aluno 1 com fotos aninhadas, obtido %+v", aluno)
	}
}

func TestObterAlunoInexistente(t *testing.T) {
	setupControllerTest(t)
	req := withID(httptest.NewRequest("GET", "/api/v1/alunos/999", nil), "999")
	rec := httptest.NewRecorder()

	ObterAluno(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("esperado 404, obtido %d", rec.Code)
	}
}

func TestCriarAluno(t *testing.T) {
	setupControllerTest(t)
	body := strings.NewReader(`{"nome":"Novo Aluno","foto":"https://example.com/foto.jpg","turma":"2026.1"}`)
	req := httptest.NewRequest("POST", "/api/v1/alunos", body)
	rec := httptest.NewRecorder()

	CriarAluno(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtido %d", rec.Code)
	}
	var aluno model.Aluno
	if err := json.NewDecoder(rec.Body).Decode(&aluno); err != nil {
		t.Fatal("resposta nao e JSON valido")
	}
	if aluno.ID == 0 || aluno.Nome != "Novo Aluno" {
		t.Fatalf("aluno criado invalido: %+v", aluno)
	}
}

func TestCriarAlunoRejeitaURLInvalida(t *testing.T) {
	setupControllerTest(t)
	body := strings.NewReader(`{"nome":"Novo Aluno","foto":"javascript:alert(1)","turma":"2026.1"}`)
	req := httptest.NewRequest("POST", "/api/v1/alunos", body)
	rec := httptest.NewRecorder()

	CriarAluno(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtido %d", rec.Code)
	}
}

func TestAtualizarAluno(t *testing.T) {
	setupControllerTest(t)
	body := strings.NewReader(`{"nome":"Atualizado","foto":"https://example.com/nova.jpg","turma":"2026.2"}`)
	req := withID(httptest.NewRequest("PUT", "/api/v1/alunos/1", body), "1")
	rec := httptest.NewRecorder()

	AtualizarAluno(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}
	var aluno model.Aluno
	if err := json.NewDecoder(rec.Body).Decode(&aluno); err != nil {
		t.Fatal("resposta nao e JSON valido")
	}
	if aluno.Nome != "Atualizado" {
		t.Fatalf("esperado nome Atualizado, obtido %q", aluno.Nome)
	}
}

func TestRemoverAluno(t *testing.T) {
	setupControllerTest(t)
	req := withID(httptest.NewRequest("DELETE", "/api/v1/alunos/2", nil), "2")
	rec := httptest.NewRecorder()

	RemoverAluno(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("esperado 204, obtido %d", rec.Code)
	}
}

func TestListarEventos(t *testing.T) {
	setupControllerTest(t)
	req := httptest.NewRequest("GET", "/api/v1/eventos", nil)
	rec := httptest.NewRecorder()

	ListarEventos(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}
	var eventos []model.Evento
	if err := json.NewDecoder(rec.Body).Decode(&eventos); err != nil {
		t.Fatal("resposta nao e JSON valido")
	}
	if len(eventos) < 2 {
		t.Fatal("esperava pelo menos 2 eventos")
	}
}

func TestCriarEvento(t *testing.T) {
	setupControllerTest(t)
	body := strings.NewReader(`{"titulo":"Demo Day","descricao":"Apresentacao do MVP","data":"2026-07-01","imagem_url":"https://example.com/demo.jpg"}`)
	req := httptest.NewRequest("POST", "/api/v1/eventos", body)
	rec := httptest.NewRecorder()

	CriarEvento(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtido %d", rec.Code)
	}
}

func TestLoginRetornaJWT(t *testing.T) {
	setupControllerTest(t)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"usuario":"admin","senha":"admin123"}`))
	rec := httptest.NewRecorder()

	Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}
	var tokens tokenResponse
	if err := json.NewDecoder(rec.Body).Decode(&tokens); err != nil {
		t.Fatal("resposta nao e JSON valido")
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" || tokens.TokenType != "Bearer" {
		t.Fatalf("tokens invalidos: %+v", tokens)
	}
}

func TestRefreshRotacionaToken(t *testing.T) {
	setupControllerTest(t)
	loginReq := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"usuario":"admin","senha":"admin123"}`))
	loginRec := httptest.NewRecorder()
	Login(loginRec, loginReq)

	var loginTokens tokenResponse
	if err := json.NewDecoder(loginRec.Body).Decode(&loginTokens); err != nil {
		t.Fatal(err)
	}

	refreshReq := httptest.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(`{"refresh_token":"`+loginTokens.RefreshToken+`"}`))
	refreshRec := httptest.NewRecorder()
	Refresh(refreshRec, refreshReq)

	if refreshRec.Code != http.StatusOK {
		t.Fatalf("esperado 200 no refresh, obtido %d", refreshRec.Code)
	}
	var rotated tokenResponse
	if err := json.NewDecoder(refreshRec.Body).Decode(&rotated); err != nil {
		t.Fatal(err)
	}
	if rotated.RefreshToken == "" || rotated.RefreshToken == loginTokens.RefreshToken {
		t.Fatal("refresh token deve ser rotacionado")
	}

	reuseReq := httptest.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(`{"refresh_token":"`+loginTokens.RefreshToken+`"}`))
	reuseRec := httptest.NewRecorder()
	Refresh(reuseRec, reuseReq)
	if reuseRec.Code != http.StatusUnauthorized {
		t.Fatalf("refresh token antigo deve ser recusado, obtido %d", reuseRec.Code)
	}
}
