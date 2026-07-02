package repository

import (
	"context"
	"time"

	"anuario/interno/model"
)

type Store interface {
	ListarAlunos(context.Context) ([]model.Aluno, error)
	BuscarAlunoPorID(context.Context, int) (model.Aluno, bool, error)
	BuscarAlunoComFotos(context.Context, int) (model.AlunoComFotos, bool, error)
	AdicionarAluno(context.Context, model.AlunoInput) (model.Aluno, error)
	AtualizarAluno(context.Context, int, model.AlunoInput) (model.Aluno, bool, error)
	RemoverAluno(context.Context, int) (bool, error)

	ListarEventos(context.Context) ([]model.Evento, error)
	BuscarEventoPorID(context.Context, int) (model.Evento, bool, error)
	CriarEvento(context.Context, model.EventoInput) (model.Evento, error)
	AtualizarEvento(context.Context, int, model.EventoInput) (model.Evento, bool, error)
	RemoverEvento(context.Context, int) (bool, error)

	ListarFotosPorAluno(context.Context, int) ([]model.Foto, error)
	CriarFoto(context.Context, int, model.FotoInput) (model.Foto, bool, error)

	BuscarAdminPorUsuario(context.Context, string) (model.Admin, bool, error)
	BuscarAdminPorID(context.Context, int) (model.Admin, bool, error)
	CriarRefreshToken(context.Context, int, string, time.Time) error
	BuscarRefreshToken(context.Context, string) (model.RefreshToken, bool, error)
	RevogarRefreshToken(context.Context, string) error
}
