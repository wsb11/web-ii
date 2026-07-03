package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"anuario/interno/auth"
	"anuario/interno/controller"
	"anuario/interno/model"
	"anuario/interno/repository"
)

func setupRouterTest(t *testing.T) (*auth.JWTManager, http.Handler) {
	t.Helper()

	manager := auth.MustJWTManager("test-secret-anuario-2026", time.Minute)
	controller.Configurar(repository.NewMemoryStore(), manager)

	return manager, router(manager)
}

func TestRouterProtegeRotasDeAlunosComJWT(t *testing.T) {
	_, handler := setupRouterTest(t)

	tests := []struct {
		name string
		path string
	}{
		{name: "lista", path: "/api/v1/alunos"},
		{name: "detalhe", path: "/api/v1/alunos/1"},
		{name: "fotos", path: "/api/v1/alunos/1/fotos"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("esperado 401 sem JWT, obtido %d", rec.Code)
			}
		})
	}
}

func TestRouterListaAlunosComJWTAdmin(t *testing.T) {
	manager, handler := setupRouterTest(t)
	token, err := manager.Generate(model.Admin{ID: 1, Usuario: "admin", Role: "admin"})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alunos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200 com JWT admin, obtido %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatal("Content-Type deve ser application/json")
	}
}
