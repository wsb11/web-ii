package controller

import (
	"time"

	"anuario/interno/auth"
	"anuario/interno/repository"
)

var (
	store      repository.Store = repository.NewMemoryStore()
	tokenMaker                  = auth.MustJWTManager("dev-secret-anuario-2026", 15*time.Minute)
	refreshTTL                  = 7 * 24 * time.Hour
)

func Configurar(deps repository.Store, manager *auth.JWTManager) {
	if deps != nil {
		store = deps
	}
	if manager != nil {
		tokenMaker = manager
	}
}

func ConfigurarRefreshTTL(ttl time.Duration) {
	if ttl > 0 {
		refreshTTL = ttl
	}
}
