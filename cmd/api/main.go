package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"anuario/interno/auth"
	"anuario/interno/controller"
	"anuario/interno/db"
	customMiddleware "anuario/interno/middleware"
	"anuario/interno/repository"
)

func main() {
	ctx := context.Background()
	appStore := repository.Store(repository.NewMemoryStore())

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		database, err := sql.Open("pgx", databaseURL)
		if err != nil {
			log.Fatal(err)
		}
		defer database.Close()

		if err := database.PingContext(ctx); err != nil {
			log.Fatal(err)
		}
		if err := db.Migrate(ctx, database); err != nil {
			log.Fatal(err)
		}

		pgStore := repository.NewPostgresStore(database)
		adminHash, err := auth.HashPassword(getenv("ADMIN_PASSWORD", "admin123"))
		if err != nil {
			log.Fatal(err)
		}
		if err := pgStore.EnsureAdmin(ctx, getenv("ADMIN_USER", "admin"), adminHash); err != nil {
			log.Fatal(err)
		}
		if err := pgStore.SeedDemoData(ctx); err != nil {
			log.Fatal(err)
		}
		appStore = pgStore
		log.Println("Persistencia PostgreSQL habilitada")
	} else {
		log.Println("DATABASE_URL ausente; usando armazenamento em memoria para desenvolvimento")
	}

	tokenMaker, err := auth.NewJWTManager(getenv("JWT_SECRET", "dev-secret-anuario-2026"), 15*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	controller.Configurar(appStore, tokenMaker)

	addr := ":" + getenv("PORT", "8080")
	log.Printf("Servidor rodando em http://localhost%s", addr)
	if err := http.ListenAndServe(addr, router(tokenMaker)); err != nil {
		log.Fatal(err)
	}
}

func router(tokenMaker *auth.JWTManager) http.Handler {
	r := chi.NewRouter()
	authLimiter := customMiddleware.NewRateLimiter(5, time.Minute)
	fileServer := http.FileServer(http.Dir("public"))

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(customMiddleware.SecurityHeaders)
	r.Use(customMiddleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.With(authLimiter.Middleware).Post("/auth/login", controller.Login)
		r.With(authLimiter.Middleware).Post("/auth/refresh", controller.Refresh)

		r.Get("/eventos", controller.ListarEventos)
		r.Get("/eventos/{id}", controller.ObterEvento)

		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.AutenticacaoJWT(tokenMaker))
			r.Use(customMiddleware.AutorizarRole("admin"))

			registrarLeituraAlunos(r, "/alunos")
			registrarLeituraAlunos(r, "/public/alunos")
			r.Post("/alunos", controller.CriarAluno)
			r.Put("/alunos/{id}", controller.AtualizarAluno)
			r.Delete("/alunos/{id}", controller.RemoverAluno)
			r.Post("/alunos/{id}/fotos", controller.CriarFotoDoAluno)

			r.Post("/eventos", controller.CriarEvento)
			r.Put("/eventos/{id}", controller.AtualizarEvento)
			r.Delete("/eventos/{id}", controller.RemoverEvento)
		})
	})

	r.Handle("/*", fileServer)

	return r
}

func registrarLeituraAlunos(r chi.Router, base string) {
	r.Get(base, controller.ListarAlunos)
	r.Get(base+"/{id}", controller.ObterAluno)
	r.Get(base+"/{id}/fotos", controller.ListarFotosDoAluno)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
