package repository

import (
	"context"
	"testing"

	"anuario/interno/model"
)

func TestMemoryStoreAdicionarAluno(t *testing.T) {
	repo := NewMemoryStore()
	alunoSalvo, err := repo.AdicionarAluno(context.Background(), model.AlunoInput{Nome: "Teste da Silva", Turma: "2026.1"})

	if err != nil {
		t.Fatal(err)
	}
	if alunoSalvo.ID == 0 {
		t.Fatalf("esperava ID valido, recebeu %d", alunoSalvo.ID)
	}
	if alunoSalvo.Nome != "Teste da Silva" {
		t.Fatalf("esperava nome Teste da Silva, recebeu %q", alunoSalvo.Nome)
	}
}

func TestMemoryStoreBuscarPorIDExistente(t *testing.T) {
	repo := NewMemoryStore()
	aluno, ok, err := repo.BuscarAlunoPorID(context.Background(), 1)

	if err != nil {
		t.Fatal(err)
	}
	if !ok || aluno.ID != 1 {
		t.Fatal("deveria encontrar o aluno ID 1")
	}
}

func TestMemoryStoreBuscarPorIDInexistente(t *testing.T) {
	repo := NewMemoryStore()
	_, ok, err := repo.BuscarAlunoPorID(context.Background(), 999)

	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("nao deveria encontrar aluno com ID 999")
	}
}

func TestMemoryStoreAlunoComFotos(t *testing.T) {
	repo := NewMemoryStore()
	aluno, ok, err := repo.BuscarAlunoComFotos(context.Background(), 1)

	if err != nil {
		t.Fatal(err)
	}
	if !ok || len(aluno.Fotos) == 0 {
		t.Fatalf("esperava aluno com fotos, obtido %+v", aluno)
	}
}
