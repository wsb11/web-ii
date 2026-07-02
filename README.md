
# Memórias — Anuário Digital
### Documentação Técnica e de Produto | UFRN / IMD | Versão 1.0 — 2026

---

## Autores

| Nome | Matrícula |
|------|-----------|
| Vinicios David Martins Bezerra | 20220062699 |
| Weuler dos Santos Barbosa | 20220026576 |

---

## Sumário

1. [Visão Geral do Produto](#1-visão-geral-do-produto)
2. [Equipe e Responsabilidades](#2-equipe-e-responsabilidades)
3. [Stack Tecnológico](#3-stack-tecnológico)
4. [Definição do MVP](#4-definição-do-mvp)
5. [Product Backlog](#5-product-backlog)
6. [Arquitetura do Sistema](#6-arquitetura-do-sistema)
7. [Modelagem do Banco de Dados](#7-modelagem-do-banco-de-dados)
8. [Fluxo da Aplicação](#8-fluxo-da-aplicação)
9. [Considerações Finais](#9-considerações-finais)

---

## 1. Visão Geral do Produto

### 1.1 Descrição

O **Memórias — Anuário Digital** é um sistema web desenvolvido para registro, organização e visualização de histórias acadêmicas de turmas universitárias. O projeto foi realizado no âmbito do Instituto Metrópole Digital (IMD) da Universidade Federal do Rio Grande do Norte (UFRN).

A plataforma permite preservar memórias, trajetórias, conquistas e momentos marcantes de uma turma ao longo dos anos, oferecendo uma alternativa digital, acessível e duradoura em relação aos tradicionais registros físicos.

### 1.2 Declaração de Visão

> **Para** alunos, professores e comunidade escolar
>
> **Que** não possuem um registro estruturado e duradouro das memórias, trajetórias e vivências das turmas,
>
> **O Memórias (Anuário Digital) é um** sistema web de registro e visualização de histórias acadêmicas
>
> **Que** permite armazenar, organizar e reviver momentos, perfis e conquistas de uma turma ao longo dos anos.
>
> **Diferente de** placas físicas e registros tradicionais limitados e temporários,
>
> **Nosso produto** centraliza tudo digitalmente com acesso contínuo, multimídia e editável.

### 1.3 Problema Central

Falta de um sistema digital estruturado para registrar e acessar memórias completas de uma turma, indo além de registros físicos limitados como placas, álbuns e documentos impressos que se deterioram com o tempo e não oferecem acesso remoto ou multimídia.

### 1.4 Hipótese de Valor

> *"Acreditamos que alunos e instituições vão usar e alimentar o anuário digital porque ele permite preservar memórias, histórias e identidade da turma de forma acessível e permanente."*

---

## 2. Equipe e Responsabilidades

### 2.1 Composição da Equipe

| Nome | Matrícula | Responsabilidade |
|------|-----------|------------------|
| Emilly Miller Moreira | 20230031417 | Front-end — Interfaces das páginas (início, alunos, galeria, eventos) |
| Francisco Matheus Fonseca De Farias | 20220052923 | Front-end — Interfaces das páginas (início, alunos, galeria, eventos) |
| Vinicios David Martins Bezerra | 20220062699 | Backend em Go — Regra de negócio, validações e comportamentos do sistema |
| Weuler dos Santos Barbosa | 20220026576 | Backend em Go — Camada de rotas e handlers, estruturação das requisições HTTP |

### 2.2 Responsabilidade Compartilhada

Todos os membros ficaram responsáveis pelo **banco de dados** e integração com **PostgreSQL via sqlc**, modelando tabelas, queries e persistência dos dados.

---

## 3. Stack Tecnológico

### 3.1 Tecnologias Utilizadas

| Camada | Tecnologia | Descrição |
|--------|------------|-----------|
| Front-end | HTML5 | Estrutura e marcação das páginas web |
| Front-end | CSS3 | Estilização e layout das interfaces |
| Front-end | JavaScript | Interatividade e dinamismo das páginas |
| Back-end | Go (Golang) | Linguagem principal do servidor; alto desempenho e tipagem estática |
| Back-end | Chi | Roteador HTTP leve para gerenciamento de rotas e middlewares |
| Banco de Dados | PostgreSQL | Sistema de gerenciamento de banco de dados relacional |
| Banco de Dados | sqlc | Geração de queries SQL tipadas, garantindo segurança de tipos |

### 3.2 Justificativas Técnicas

#### Go + Chi (Back-end)

Go foi escolhido por sua performance, simplicidade e modelo de concorrência eficiente, sendo ideal para APIs web. O roteador Chi é leve, compatível com a biblioteca padrão `net/http` e oferece suporte a middlewares de forma modular.

#### PostgreSQL + sqlc

PostgreSQL oferece robustez, suporte a transações ACID e escalabilidade. O sqlc permite escrever queries SQL puras e gerar automaticamente código Go tipado, eliminando erros de runtime relacionados ao banco de dados.

#### Front-end sem frameworks

A opção por HTML, CSS e JavaScript puro reduz a complexidade de build e dependências externas, facilitando a manutenção e o aprendizado dos membros responsáveis pelo front-end.

---

## 4. Definição do MVP

### 4.1 Funcionalidades Essenciais

#### P1 — Core (Obrigatório)

- Cadastro e listagem de alunos
- Página individual do aluno (perfil básico)
- Galeria de fotos da turma
- Página de eventos (linha do tempo)
- Visualização pública (visitante sem autenticação)

#### P1 — Admin (Obrigatório)

- Login de administrador
- CRUD de alunos (adicionar, editar, excluir)
- Upload de fotos
- Cadastro de eventos

#### P2 — Importante (Próximas Sprints)

- Perfis de professores
- Edição limitada por aluno
- Comentários e opiniões

### 4.2 Fora do Escopo do MVP

- Aplicativo mobile
- Sistema avançado de permissões
- Interações sociais (curtidas, comentários complexos)
- Integrações com serviços externos

### 4.3 Critérios de Sucesso

1. Usuário consegue acessar e navegar pelo anuário sem erros.
2. Administrador consegue cadastrar alunos, fotos e eventos.
3. Sistema armazena e exibe corretamente os dados persistidos.
4. Experiência básica funcional, intuitiva e responsiva.

---

## 5. Product Backlog

### 5.1 Planejamento de Sprints

| Prioridade | User Story | Critérios de Aceitação | Sprint |
|------------|------------|------------------------|--------|
| P1 | Como visitante, quero visualizar o anuário | Página inicial + navegação funcional | 1 |
| P1 | Como visitante, quero ver alunos da turma | Lista com nome + foto | 1 |
| P1 | Como visitante, quero ver detalhes de um aluno | Modal/página com informações básicas | 2 |
| P1 | Como visitante, quero ver galeria de fotos | Grid de imagens funcional | 2 |
| P1 | Como visitante, quero ver eventos em linha do tempo | Ordenação cronológica | 2 |
| P1 | Como admin, quero fazer login | Autenticação funcionando | 3 |
| P1 | Como admin, quero cadastrar alunos | Formulário com validação | 3 |
| P1 | Como admin, quero editar/excluir alunos | CRUD completo | 3 |
| P1 | Como admin, quero cadastrar fotos | Upload + exibição | 4 |
| P1 | Como admin, quero cadastrar eventos | Texto + data + imagem | 4 |
| P2 | Como professor, quero editar meu perfil | Campos editáveis básicos | 5 |
| P2 | Como aluno, quero editar meu perfil | Permissão limitada | 5 |
| P3 | Como usuário, quero comentar | Sistema simples de comentários | 6 |

---

## 6. Arquitetura do Sistema

### 6.1 Visão Geral da Arquitetura

O sistema segue uma arquitetura em camadas (*layered architecture*), com separação clara entre front-end, back-end e persistência de dados. A comunicação entre o cliente e o servidor se dá via HTTP, com o servidor Go respondendo a requisições REST.

| Camada | Descrição |
|--------|-----------|
| Apresentação | HTML, CSS, JavaScript — renderiza as páginas e consome a API REST |
| Rotas / Handlers | Go + Chi — recebe requisições HTTP, valida entrada e retorna respostas |
| Regra de Negócio | Go — validações, lógica de domínio e orquestração de operações |
| Acesso a Dados | sqlc — queries tipadas geradas a partir de SQL puro |
| Persistência | PostgreSQL — armazenamento relacional de alunos, eventos e fotos |

### 6.2 Diagrama de Camadas

```
┌─────────────────────────────────────────────────┐
│              Cliente (Browser)                  │
│          HTML  •  CSS  •  JavaScript            │
└─────────────────────┬───────────────────────────┘
                      │  HTTP / REST
┌─────────────────────▼───────────────────────────┐
│              Servidor Go                        │
│  ┌──────────────────────────────────────────┐   │
│  │  Rotas & Handlers (Chi)                  │   │
│  ├──────────────────────────────────────────┤   │
│  │  Regra de Negócio                        │   │
│  ├──────────────────────────────────────────┤   │
│  │  Acesso a Dados (sqlc)                   │   │
│  └──────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────┘
                      │  SQL
┌─────────────────────▼───────────────────────────┐
│              PostgreSQL                         │
│       alunos  •  eventos  •  fotos              │
└─────────────────────────────────────────────────┘
```

### 6.3 Módulos Principais

#### Módulo de Alunos
Responsável pelo cadastro, edição, exclusão e listagem de alunos. Cada aluno possui perfil básico com nome, foto, curso e informações complementares.

#### Módulo de Galeria
Gerencia o upload, armazenamento e exibição de fotos da turma. As imagens são organizadas em grade e associadas à turma correspondente.

#### Módulo de Eventos
Registra e exibe eventos em ordem cronológica. Cada evento contém título, data, descrição textual e imagem opcional, formando uma linha do tempo da turma.

#### Módulo de Autenticação
Implementa o login do administrador com validação de credenciais. O administrador possui permissão total de CRUD sobre os demais módulos.

---

## 7. Modelagem do Banco de Dados

### 7.1 Entidades Principais

#### Tabela: `alunos`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | `SERIAL` / `UUID` | Identificador único do aluno |
| `nome` | `VARCHAR(255)` | Nome completo do aluno |
| `matricula` | `VARCHAR(20)` | Número de matrícula na UFRN |
| `foto_url` | `TEXT` | URL da foto de perfil |
| `descricao` | `TEXT` | Texto de apresentação / bio |
| `criado_em` | `TIMESTAMP` | Data de cadastro no sistema |

#### Tabela: `eventos`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | `SERIAL` / `UUID` | Identificador único do evento |
| `titulo` | `VARCHAR(255)` | Título do evento |
| `descricao` | `TEXT` | Descrição detalhada do evento |
| `data_evento` | `DATE` | Data em que o evento ocorreu |
| `imagem_url` | `TEXT` | URL da imagem do evento (opcional) |

#### Tabela: `fotos`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | `SERIAL` / `UUID` | Identificador único da foto |
| `url` | `TEXT` | Caminho de acesso à imagem |
| `legenda` | `VARCHAR(255)` | Legenda descritiva da foto (opcional) |
| `enviado_em` | `TIMESTAMP` | Data e hora do upload |

#### Tabela: `admins`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | `SERIAL` | Identificador único do administrador |
| `usuario` | `VARCHAR(100)` | Nome de usuário para login |
| `senha_hash` | `TEXT` | Hash da senha (bcrypt) |
| `criado_em` | `TIMESTAMP` | Data de criação da conta |

### 7.2 Exemplo de Query sqlc

```sql
-- name: ListarAlunos :many
SELECT id, nome, matricula, foto_url, descricao
FROM alunos
ORDER BY nome ASC;

-- name: BuscarAlunoPorID :one
SELECT id, nome, matricula, foto_url, descricao
FROM alunos
WHERE id = $1;

-- name: CriarAluno :one
INSERT INTO alunos (nome, matricula, foto_url, descricao)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: AtualizarAluno :one
UPDATE alunos
SET nome = $2, matricula = $3, foto_url = $4, descricao = $5
WHERE id = $1
RETURNING *;

-- name: DeletarAluno :exec
DELETE FROM alunos WHERE id = $1;
```

### 7.3 Integração com sqlc

O sqlc é utilizado para gerar automaticamente structs Go e funções de acesso ao banco a partir de arquivos `.sql`. Esse fluxo garante que qualquer alteração no schema reflita erros de compilação no código, prevenindo inconsistências em tempo de execução.

---

## 8. Fluxo da Aplicação

### 8.1 Fluxo do Visitante

```
Página Inicial
     │
     ├──► Lista de Alunos ──► Perfil Individual do Aluno
     │
     ├──► Galeria de Fotos
     │
     └──► Linha do Tempo de Eventos
```

1. Acessa a página inicial pública do anuário.
2. Navega pela lista de alunos e visualiza perfis individuais.
3. Acessa a galeria de fotos da turma.
4. Consulta os eventos na linha do tempo cronológica.

### 8.2 Fluxo do Administrador

```
Login
  │
  └──► Painel Admin
            │
            ├──► Gerenciar Alunos (CRUD)
            │         ├── Criar
            │         ├── Editar
            │         └── Excluir
            │
            ├──► Upload de Fotos
            │
            └──► Cadastrar Eventos
```

1. Acessa a rota de login e autentica-se com credenciais.
2. Gerencia alunos: cria, edita dados e exclui registros (CRUD completo).
3. Realiza upload de fotos para a galeria.
4. Cadastra eventos com título, descrição, data e imagem.
5. Encerra sessão ao sair do painel.

### 8.3 Rotas da API

| Método | Rota | Descrição | Autenticação |
|--------|------|-----------|:------------:|
| `GET` | `/` | Página inicial pública | — |
| `GET` | `/alunos` | Listagem de alunos | — |
| `GET` | `/alunos/{id}` | Perfil de um aluno | — |
| `GET` | `/galeria` | Galeria de fotos | — |
| `GET` | `/eventos` | Linha do tempo | — |
| `POST` | `/admin/login` | Autenticação do admin | — |
| `POST` | `/admin/alunos` | Criar aluno | ✅ |
| `PUT` | `/admin/alunos/{id}` | Editar aluno | ✅ |
| `DELETE` | `/admin/alunos/{id}` | Excluir aluno | ✅ |
| `POST` | `/admin/fotos` | Upload de foto | ✅ |
| `POST` | `/admin/eventos` | Cadastrar evento | ✅ |

---

## 9. Considerações Finais

### 9.1 Desafios e Aprendizados

O desenvolvimento do Memórias proporcionou à equipe experiência prática com a linguagem Go em ambiente web, o paradigma de queries tipadas com sqlc e a integração de um back-end estruturado com front-end em HTML/CSS/JavaScript puro. A divisão clara de responsabilidades favoreceu a colaboração e permitiu entregas incrementais por sprint.

### 9.2 Evolução Futura

Com a base estabelecida pelo MVP, as próximas iterações poderão incorporar os itens do backlog de prioridade P2 e P3, como perfis de professores, edição de perfil pelo próprio aluno e um sistema de comentários. A longo prazo, o produto pode evoluir para suportar múltiplas turmas e instituições distintas.

### 9.3 Conclusão

O projeto Memórias demonstra como tecnologias modernas e de código aberto — Go, PostgreSQL e sqlc — podem ser combinadas para construir uma aplicação web funcional, escalável e de baixo custo de manutenção. O anuário digital representa uma solução real para a preservação da identidade e das memórias de turmas universitárias, contribuindo para a história do Instituto Metrópole Digital.

---

*Universidade Federal do Rio Grande do Norte  •  Instituto Metrópole Digital  •  2026*
