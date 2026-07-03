package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"anuario/interno/auth"
	"anuario/interno/model"
)

type MemoryStore struct {
	mu           sync.RWMutex
	alunos       []model.Aluno
	eventos      []model.Evento
	fotos        []model.Foto
	admins       []model.Admin
	refresh      []model.RefreshToken
	nextAlunoID  int
	nextEventoID int
	nextFotoID   int
	nextTokenID  int
}

func NewMemoryStore() *MemoryStore {
	alunos := demoAlunosComID()
	fotos := demoFotos()

	return &MemoryStore{
		alunos: alunos,
		eventos: []model.Evento{
			{ID: 1, Titulo: "Aula Inaugural", Descricao: "Primeira aula do semestre", Data: "2026-02-15"},
			{ID: 2, Titulo: "Semana de Tecnologia", Descricao: "Palestras e workshops", Data: "2026-06-20"},
			{ID: 3, Titulo: "Formatura", Descricao: "Colacao de grau", Data: "2029-12-10"},
			{ID: 4, Titulo: "Feira de Estagio", Descricao: "Oportunidades de estagio", Data: "2026-09-05"},
		},
		fotos: fotos,
		admins: []model.Admin{
			{ID: 1, Usuario: "admin", SenhaHash: auth.MustHashPassword("admin123"), Role: "admin"},
		},
		nextAlunoID:  len(alunos) + 1,
		nextEventoID: 5,
		nextFotoID:   len(fotos) + 1,
		nextTokenID:  1,
	}
}

func (s *MemoryStore) ListarAlunos(context.Context) ([]model.Aluno, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Aluno(nil), s.alunos...), nil
}

func (s *MemoryStore) BuscarAlunoPorID(_ context.Context, id int) (model.Aluno, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, aluno := range s.alunos {
		if aluno.ID == id {
			return aluno, true, nil
		}
	}
	return model.Aluno{}, false, nil
}

func (s *MemoryStore) BuscarAlunoComFotos(_ context.Context, id int) (model.AlunoComFotos, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var aluno model.Aluno
	encontrado := false
	for _, item := range s.alunos {
		if item.ID == id {
			aluno = item
			encontrado = true
			break
		}
	}
	if !encontrado {
		return model.AlunoComFotos{}, false, nil
	}

	fotos := make([]model.Foto, 0)
	for _, foto := range s.fotos {
		if foto.AlunoID == id {
			fotos = append(fotos, foto)
		}
	}
	return model.AlunoComFotos{Aluno: aluno, Fotos: fotos}, true, nil
}

func (s *MemoryStore) AdicionarAluno(_ context.Context, input model.AlunoInput) (model.Aluno, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	aluno := model.Aluno{ID: s.nextAlunoID, Nome: input.Nome, Foto: input.Foto, Turma: input.Turma}
	s.nextAlunoID++
	s.alunos = append(s.alunos, aluno)
	return aluno, nil
}

func (s *MemoryStore) AtualizarAluno(_ context.Context, id int, input model.AlunoInput) (model.Aluno, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, aluno := range s.alunos {
		if aluno.ID == id {
			atualizado := model.Aluno{ID: id, Nome: input.Nome, Foto: input.Foto, Turma: input.Turma}
			s.alunos[i] = atualizado
			return atualizado, true, nil
		}
	}
	return model.Aluno{}, false, nil
}

func (s *MemoryStore) RemoverAluno(_ context.Context, id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, aluno := range s.alunos {
		if aluno.ID == id {
			s.alunos = append(s.alunos[:i], s.alunos[i+1:]...)
			fotos := s.fotos[:0]
			for _, foto := range s.fotos {
				if foto.AlunoID != id {
					fotos = append(fotos, foto)
				}
			}
			s.fotos = fotos
			return true, nil
		}
	}
	return false, nil
}

func (s *MemoryStore) ListarEventos(context.Context) ([]model.Evento, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	eventos := append([]model.Evento(nil), s.eventos...)
	sort.Slice(eventos, func(i, j int) bool {
		return eventos[i].Data < eventos[j].Data
	})
	return eventos, nil
}

func (s *MemoryStore) BuscarEventoPorID(_ context.Context, id int) (model.Evento, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, evento := range s.eventos {
		if evento.ID == id {
			return evento, true, nil
		}
	}
	return model.Evento{}, false, nil
}

func (s *MemoryStore) CriarEvento(_ context.Context, input model.EventoInput) (model.Evento, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	evento := model.Evento{
		ID:        s.nextEventoID,
		Titulo:    input.Titulo,
		Descricao: input.Descricao,
		Data:      input.Data,
		ImagemURL: input.ImagemURL,
	}
	s.nextEventoID++
	s.eventos = append(s.eventos, evento)
	return evento, nil
}

func (s *MemoryStore) AtualizarEvento(_ context.Context, id int, input model.EventoInput) (model.Evento, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, evento := range s.eventos {
		if evento.ID == id {
			atualizado := model.Evento{ID: id, Titulo: input.Titulo, Descricao: input.Descricao, Data: input.Data, ImagemURL: input.ImagemURL}
			s.eventos[i] = atualizado
			return atualizado, true, nil
		}
	}
	return model.Evento{}, false, nil
}

func (s *MemoryStore) RemoverEvento(_ context.Context, id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, evento := range s.eventos {
		if evento.ID == id {
			s.eventos = append(s.eventos[:i], s.eventos[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

func (s *MemoryStore) ListarFotosPorAluno(_ context.Context, alunoID int) ([]model.Foto, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fotos := make([]model.Foto, 0)
	for _, foto := range s.fotos {
		if foto.AlunoID == alunoID {
			fotos = append(fotos, foto)
		}
	}
	return fotos, nil
}

func (s *MemoryStore) CriarFoto(_ context.Context, alunoID int, input model.FotoInput) (model.Foto, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existe := false
	for _, aluno := range s.alunos {
		if aluno.ID == alunoID {
			existe = true
			break
		}
	}
	if !existe {
		return model.Foto{}, false, nil
	}

	foto := model.Foto{ID: s.nextFotoID, AlunoID: alunoID, URL: input.URL, Legenda: input.Legenda}
	s.nextFotoID++
	s.fotos = append(s.fotos, foto)
	return foto, true, nil
}

func (s *MemoryStore) BuscarAdminPorUsuario(_ context.Context, usuario string) (model.Admin, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, admin := range s.admins {
		if admin.Usuario == usuario {
			return admin, true, nil
		}
	}
	return model.Admin{}, false, nil
}

func (s *MemoryStore) BuscarAdminPorID(_ context.Context, id int) (model.Admin, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, admin := range s.admins {
		if admin.ID == id {
			return admin, true, nil
		}
	}
	return model.Admin{}, false, nil
}

func (s *MemoryStore) CriarRefreshToken(_ context.Context, adminID int, tokenHash string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	token := model.RefreshToken{ID: s.nextTokenID, AdminID: adminID, TokenHash: tokenHash, ExpiresAt: expiresAt}
	s.nextTokenID++
	s.refresh = append(s.refresh, token)
	return nil
}

func (s *MemoryStore) BuscarRefreshToken(_ context.Context, tokenHash string) (model.RefreshToken, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	for _, token := range s.refresh {
		if token.TokenHash == tokenHash && token.RevokedAt == nil && token.ExpiresAt.After(now) {
			return token, true, nil
		}
	}
	return model.RefreshToken{}, false, nil
}

func (s *MemoryStore) RevogarRefreshToken(_ context.Context, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for i, token := range s.refresh {
		if token.TokenHash == tokenHash && token.RevokedAt == nil {
			s.refresh[i].RevokedAt = &now
			return nil
		}
	}
	return nil
}
