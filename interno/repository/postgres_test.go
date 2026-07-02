package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"anuario/interno/db"
	"anuario/interno/model"
)

func TestPostgresStoreCRUDComRelacionamento(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL nao configurado")
	}

	ctx := context.Background()
	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(ctx, database); err != nil {
		t.Fatal(err)
	}
	if _, err := database.ExecContext(ctx, "TRUNCATE refresh_tokens, fotos, eventos, alunos, admins RESTART IDENTITY CASCADE"); err != nil {
		t.Fatal(err)
	}

	repo := NewPostgresStore(database)
	aluno, err := repo.AdicionarAluno(ctx, model.AlunoInput{Nome: "Aluno Postgres", Foto: "https://example.com/aluno.jpg", Turma: "2026.1"})
	if err != nil {
		t.Fatal(err)
	}
	if _, found, err := repo.CriarFoto(ctx, aluno.ID, model.FotoInput{URL: "https://example.com/foto.jpg", Legenda: "Foto do aluno"}); err != nil || !found {
		t.Fatalf("esperava criar foto vinculada, found=%v err=%v", found, err)
	}

	detalhe, found, err := repo.BuscarAlunoComFotos(ctx, aluno.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !found || len(detalhe.Fotos) != 1 {
		t.Fatalf("esperava relacionamento 1:N carregado, obtido %+v", detalhe)
	}

	evento, err := repo.CriarEvento(ctx, model.EventoInput{Titulo: "Demo Day", Descricao: "Apresentacao do MVP", Data: "2026-07-01"})
	if err != nil {
		t.Fatal(err)
	}
	if evento.ID == 0 {
		t.Fatal("evento deveria receber ID")
	}
}
