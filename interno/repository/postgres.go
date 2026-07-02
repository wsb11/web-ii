package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"anuario/interno/db"
	"anuario/interno/model"
)

type PostgresStore struct {
	q *db.Queries
}

func NewPostgresStore(database *sql.DB) *PostgresStore {
	return &PostgresStore{q: db.New(database)}
}

func (s *PostgresStore) EnsureAdmin(ctx context.Context, usuario, senhaHash string) error {
	if _, found, err := s.BuscarAdminPorUsuario(ctx, usuario); err != nil || found {
		return err
	}
	_, err := s.q.CriarAdmin(ctx, db.CriarAdminParams{
		Usuario:   usuario,
		SenhaHash: senhaHash,
		Role:      "admin",
	})
	return err
}

func (s *PostgresStore) SeedDemoData(ctx context.Context) error {
	count, err := s.q.ContarAlunos(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	aluno, err := s.AdicionarAluno(ctx, model.AlunoInput{Nome: "Emily Miller", Foto: "https://via.placeholder.com/150", Turma: "2026.1"})
	if err != nil {
		return err
	}
	if _, err := s.AdicionarAluno(ctx, model.AlunoInput{Nome: "Francisco Matheus", Foto: "https://via.placeholder.com/150", Turma: "2026.1"}); err != nil {
		return err
	}
	if _, err := s.AdicionarAluno(ctx, model.AlunoInput{Nome: "Vinicios David", Foto: "https://via.placeholder.com/150", Turma: "2026.1"}); err != nil {
		return err
	}
	if _, _, err := s.CriarFoto(ctx, aluno.ID, model.FotoInput{URL: "https://via.placeholder.com/800x600", Legenda: "Projeto integrador"}); err != nil {
		return err
	}
	if _, err := s.CriarEvento(ctx, model.EventoInput{Titulo: "Aula Inaugural", Descricao: "Primeira aula do semestre", Data: "2026-02-15"}); err != nil {
		return err
	}
	_, err = s.CriarEvento(ctx, model.EventoInput{Titulo: "Semana de Tecnologia", Descricao: "Palestras e workshops", Data: "2026-06-20"})
	return err
}

func (s *PostgresStore) ListarAlunos(ctx context.Context) ([]model.Aluno, error) {
	rows, err := s.q.ListarAlunos(ctx)
	if err != nil {
		return nil, err
	}
	alunos := make([]model.Aluno, 0, len(rows))
	for _, row := range rows {
		alunos = append(alunos, alunoFromDB(row))
	}
	return alunos, nil
}

func (s *PostgresStore) BuscarAlunoPorID(ctx context.Context, id int) (model.Aluno, bool, error) {
	row, err := s.q.BuscarAlunoPorID(ctx, int32(id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Aluno{}, false, nil
	}
	if err != nil {
		return model.Aluno{}, false, err
	}
	return alunoFromDB(row), true, nil
}

func (s *PostgresStore) BuscarAlunoComFotos(ctx context.Context, id int) (model.AlunoComFotos, bool, error) {
	aluno, found, err := s.BuscarAlunoPorID(ctx, id)
	if err != nil || !found {
		return model.AlunoComFotos{}, found, err
	}
	fotos, err := s.ListarFotosPorAluno(ctx, id)
	if err != nil {
		return model.AlunoComFotos{}, false, err
	}
	return model.AlunoComFotos{Aluno: aluno, Fotos: fotos}, true, nil
}

func (s *PostgresStore) AdicionarAluno(ctx context.Context, input model.AlunoInput) (model.Aluno, error) {
	row, err := s.q.CriarAluno(ctx, db.CriarAlunoParams{
		Nome:  input.Nome,
		Foto:  nullString(input.Foto),
		Turma: nullString(input.Turma),
	})
	if err != nil {
		return model.Aluno{}, err
	}
	return alunoFromDB(row), nil
}

func (s *PostgresStore) AtualizarAluno(ctx context.Context, id int, input model.AlunoInput) (model.Aluno, bool, error) {
	row, err := s.q.AtualizarAluno(ctx, db.AtualizarAlunoParams{
		ID:    int32(id),
		Nome:  input.Nome,
		Foto:  nullString(input.Foto),
		Turma: nullString(input.Turma),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return model.Aluno{}, false, nil
	}
	if err != nil {
		return model.Aluno{}, false, err
	}
	return alunoFromDB(row), true, nil
}

func (s *PostgresStore) RemoverAluno(ctx context.Context, id int) (bool, error) {
	_, err := s.q.RemoverAluno(ctx, int32(id))
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (s *PostgresStore) ListarEventos(ctx context.Context) ([]model.Evento, error) {
	rows, err := s.q.ListarEventos(ctx)
	if err != nil {
		return nil, err
	}
	eventos := make([]model.Evento, 0, len(rows))
	for _, row := range rows {
		eventos = append(eventos, eventoFromDB(row))
	}
	return eventos, nil
}

func (s *PostgresStore) BuscarEventoPorID(ctx context.Context, id int) (model.Evento, bool, error) {
	row, err := s.q.BuscarEventoPorID(ctx, int32(id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Evento{}, false, nil
	}
	if err != nil {
		return model.Evento{}, false, err
	}
	return eventoFromDB(row), true, nil
}

func (s *PostgresStore) CriarEvento(ctx context.Context, input model.EventoInput) (model.Evento, error) {
	data, _ := time.Parse("2006-01-02", input.Data)
	row, err := s.q.CriarEvento(ctx, db.CriarEventoParams{
		Titulo:    input.Titulo,
		Descricao: input.Descricao,
		Data:      data,
		ImagemUrl: nullString(input.ImagemURL),
	})
	if err != nil {
		return model.Evento{}, err
	}
	return eventoFromDB(row), nil
}

func (s *PostgresStore) AtualizarEvento(ctx context.Context, id int, input model.EventoInput) (model.Evento, bool, error) {
	data, _ := time.Parse("2006-01-02", input.Data)
	row, err := s.q.AtualizarEvento(ctx, db.AtualizarEventoParams{
		ID:        int32(id),
		Titulo:    input.Titulo,
		Descricao: input.Descricao,
		Data:      data,
		ImagemUrl: nullString(input.ImagemURL),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return model.Evento{}, false, nil
	}
	if err != nil {
		return model.Evento{}, false, err
	}
	return eventoFromDB(row), true, nil
}

func (s *PostgresStore) RemoverEvento(ctx context.Context, id int) (bool, error) {
	_, err := s.q.RemoverEvento(ctx, int32(id))
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (s *PostgresStore) ListarFotosPorAluno(ctx context.Context, alunoID int) ([]model.Foto, error) {
	rows, err := s.q.ListarFotosPorAluno(ctx, int32(alunoID))
	if err != nil {
		return nil, err
	}
	fotos := make([]model.Foto, 0, len(rows))
	for _, row := range rows {
		fotos = append(fotos, fotoFromDB(row))
	}
	return fotos, nil
}

func (s *PostgresStore) CriarFoto(ctx context.Context, alunoID int, input model.FotoInput) (model.Foto, bool, error) {
	if _, found, err := s.BuscarAlunoPorID(ctx, alunoID); err != nil || !found {
		return model.Foto{}, found, err
	}
	row, err := s.q.CriarFoto(ctx, db.CriarFotoParams{
		AlunoID: int32(alunoID),
		Url:     input.URL,
		Legenda: nullString(input.Legenda),
	})
	if err != nil {
		return model.Foto{}, false, err
	}
	return fotoFromDB(row), true, nil
}

func (s *PostgresStore) BuscarAdminPorUsuario(ctx context.Context, usuario string) (model.Admin, bool, error) {
	row, err := s.q.BuscarAdminPorUsuario(ctx, usuario)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Admin{}, false, nil
	}
	if err != nil {
		return model.Admin{}, false, err
	}
	return adminFromDB(row), true, nil
}

func (s *PostgresStore) BuscarAdminPorID(ctx context.Context, id int) (model.Admin, bool, error) {
	row, err := s.q.BuscarAdminPorID(ctx, int32(id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Admin{}, false, nil
	}
	if err != nil {
		return model.Admin{}, false, err
	}
	return adminFromDB(row), true, nil
}

func (s *PostgresStore) CriarRefreshToken(ctx context.Context, adminID int, tokenHash string, expiresAt time.Time) error {
	return s.q.CriarRefreshToken(ctx, db.CriarRefreshTokenParams{
		AdminID:   int32(adminID),
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
}

func (s *PostgresStore) BuscarRefreshToken(ctx context.Context, tokenHash string) (model.RefreshToken, bool, error) {
	row, err := s.q.BuscarRefreshTokenAtivo(ctx, tokenHash)
	if errors.Is(err, sql.ErrNoRows) {
		return model.RefreshToken{}, false, nil
	}
	if err != nil {
		return model.RefreshToken{}, false, err
	}
	return refreshFromDB(row), true, nil
}

func (s *PostgresStore) RevogarRefreshToken(ctx context.Context, tokenHash string) error {
	return s.q.RevogarRefreshToken(ctx, tokenHash)
}

func alunoFromDB(row db.Aluno) model.Aluno {
	return model.Aluno{
		ID:    int(row.ID),
		Nome:  row.Nome,
		Foto:  stringFromNull(row.Foto),
		Turma: stringFromNull(row.Turma),
	}
}

func eventoFromDB(row db.Evento) model.Evento {
	return model.Evento{
		ID:        int(row.ID),
		Titulo:    row.Titulo,
		Descricao: row.Descricao,
		Data:      row.Data.Format("2006-01-02"),
		ImagemURL: stringFromNull(row.ImagemUrl),
	}
}

func fotoFromDB(row db.Foto) model.Foto {
	return model.Foto{
		ID:      int(row.ID),
		AlunoID: int(row.AlunoID),
		URL:     row.Url,
		Legenda: stringFromNull(row.Legenda),
	}
}

func adminFromDB(row db.Admin) model.Admin {
	return model.Admin{
		ID:        int(row.ID),
		Usuario:   row.Usuario,
		SenhaHash: row.SenhaHash,
		Role:      row.Role,
	}
}

func refreshFromDB(row db.RefreshToken) model.RefreshToken {
	var revokedAt *time.Time
	if row.RevokedAt.Valid {
		revokedAt = &row.RevokedAt.Time
	}
	return model.RefreshToken{
		ID:        int(row.ID),
		AdminID:   int(row.AdminID),
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		RevokedAt: revokedAt,
	}
}

func nullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

func stringFromNull(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
