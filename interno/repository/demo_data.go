package repository

import "anuario/interno/model"

func demoAlunos() []model.AlunoInput {
	return []model.AlunoInput{
		{Nome: "Ana Beatriz Bezerra Soares", Foto: "/uploads/alunos/ana-beatriz-bezerra-soares.jpg", Turma: "2026.1"},
		{Nome: "Anna Karollyne Cassiano", Foto: "/uploads/alunos/anna-karollyne-cassiano.jpg", Turma: "2026.1"},
		{Nome: "Davi Gabriel Souza de Oliveira", Foto: "/uploads/alunos/davi-gabriel-oliveira.jpg", Turma: "2026.1"},
		{Nome: "Deborah Ruth da Silva", Foto: "/uploads/alunos/deborah-ruth-silva.jpg", Turma: "2026.1"},
		{Nome: "Emanuel Kywal Pinto Cabral Filho", Foto: "/uploads/alunos/emanuel-kywal-cabral.jpg", Turma: "2026.1"},
		{Nome: "Emmanoel Pedro Fonseca de Alcantara", Foto: "/uploads/alunos/emmanoel-pedro-alcantara.jpg", Turma: "2026.1"},
		{Nome: "Felipe Matheus da Silva", Foto: "/uploads/alunos/felipe-matheus-silva.jpg", Turma: "2026.1"},
		{Nome: "Flexsivone Bezerra Oliveira", Foto: "/uploads/alunos/flexsivone-oliveira.jpg", Turma: "2026.1"},
		{Nome: "Iorrannes Firmino da Silva", Foto: "/uploads/alunos/iorrannes-firmino-silva.jpg", Turma: "2026.1"},
		{Nome: "Joao Paulo de Oliveira Cabral", Foto: "/uploads/alunos/joao-paulo-cabral.jpg", Turma: "2026.1"},
		{Nome: "Joao Pedro Pereira Frutuoso", Foto: "/uploads/alunos/joao-pedro-frutuoso.jpg", Turma: "2026.1"},
		{Nome: "Keven Diego da Rocha Barbosa", Foto: "/uploads/alunos/keven-diego-barbosa.jpg", Turma: "2026.1"},
		{Nome: "Layza Wanessa de Souza Araujo", Foto: "/uploads/alunos/layza-wanessa-araujo.jpg", Turma: "2026.1"},
		{Nome: "Leticia Gondim Guilherme", Foto: "/uploads/alunos/leticia-gondim-guilherme.jpg", Turma: "2026.1"},
		{Nome: "Lourival Cirilo de Assis Neto", Foto: "/uploads/alunos/lourival-cirilo-neto.jpg", Turma: "2026.1"},
		{Nome: "Lucas Jordan Costa da Silva", Foto: "/uploads/alunos/lucas-jordan-silva.jpg", Turma: "2026.1"},
		{Nome: "Lucas Marley de Souza Lima", Foto: "/uploads/alunos/lucas-marley-lima.jpg", Turma: "2026.1"},
		{Nome: "Luis Eduardo Pires dos Santos", Foto: "/uploads/alunos/luis-eduardo-santos.jpg", Turma: "2026.1"},
		{Nome: "Maria Joaquina Matias da Silva Oliveira", Foto: "/uploads/alunos/maria-joaquina-oliveira.jpg", Turma: "2026.1"},
		{Nome: "Maria Luiza dos Santos Silva", Foto: "/uploads/alunos/maria-luiza-silva.jpg", Turma: "2026.1"},
		{Nome: "Maria Luiza Sousa dos Santos", Foto: "/uploads/alunos/maria-luiza-sousa-santos.jpg", Turma: "2026.1"},
	}
}

func demoAlunosComID() []model.Aluno {
	inputs := demoAlunos()
	alunos := make([]model.Aluno, 0, len(inputs))
	for i, input := range inputs {
		alunos = append(alunos, model.Aluno{
			ID:    i + 1,
			Nome:  input.Nome,
			Foto:  input.Foto,
			Turma: input.Turma,
		})
	}
	return alunos
}

func demoFotos() []model.Foto {
	return []model.Foto{
		{ID: 1, AlunoID: 1, URL: "/uploads/alunos/ana-beatriz-bezerra-soares.jpg", Legenda: "Foto de perfil"},
		{ID: 2, AlunoID: 11, URL: "/uploads/alunos/joao-pedro-frutuoso.jpg", Legenda: "Foto de perfil"},
	}
}
