package auth

import (
	"strings"
	"testing"
	"time"

	"anuario/interno/model"
)

func TestHashEVerificacaoDeSenha(t *testing.T) {
	hash, err := HashPassword("senha-segura")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword("senha-segura", hash) {
		t.Fatal("senha correta deveria ser aceita")
	}
	if VerifyPassword("senha-errada", hash) {
		t.Fatal("senha incorreta nao deveria ser aceita")
	}
}

func TestJWTGenerateValidate(t *testing.T) {
	manager := MustJWTManager("test-secret-anuario-2026", time.Minute)
	token, err := manager.Generate(model.Admin{ID: 1, Usuario: "admin", Role: "admin"})
	if err != nil {
		t.Fatal(err)
	}

	claims, err := manager.Validate(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != 1 || claims.Role != "admin" {
		t.Fatalf("claims inesperadas: %+v", claims)
	}
}

func TestJWTRecusaAssinaturaAlterada(t *testing.T) {
	manager := MustJWTManager("test-secret-anuario-2026", time.Minute)
	token, err := manager.Generate(model.Admin{ID: 1, Usuario: "admin", Role: "admin"})
	if err != nil {
		t.Fatal(err)
	}

	replacement := "a"
	if strings.HasSuffix(token, "a") {
		replacement = "b"
	}
	tampered := token[:len(token)-1] + replacement
	if _, err := manager.Validate(tampered); err == nil {
		t.Fatal("token adulterado deveria ser recusado")
	}
}

func TestRefreshTokenHashEstavel(t *testing.T) {
	raw, hash, err := NewRefreshToken()
	if err != nil {
		t.Fatal(err)
	}
	if raw == "" || hash == "" {
		t.Fatal("refresh token e hash devem ser preenchidos")
	}
	if HashRefreshToken(raw) != hash {
		t.Fatal("hash do refresh token deve ser estavel")
	}
}
