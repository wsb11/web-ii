-- name: ListarAlunos :many
SELECT id, nome, foto, turma, criado_em
FROM alunos
ORDER BY id;

-- name: BuscarAlunoPorID :one
SELECT id, nome, foto, turma, criado_em
FROM alunos
WHERE id = $1;

-- name: CriarAluno :one
INSERT INTO alunos (nome, foto, turma)
VALUES ($1, $2, $3)
RETURNING id, nome, foto, turma, criado_em;

-- name: AtualizarAluno :one
UPDATE alunos
SET nome = $2, foto = $3, turma = $4
WHERE id = $1
RETURNING id, nome, foto, turma, criado_em;

-- name: RemoverAluno :one
DELETE FROM alunos
WHERE id = $1
RETURNING id;

-- name: ContarAlunos :one
SELECT COUNT(*) FROM alunos;

-- name: ListarEventos :many
SELECT id, titulo, descricao, data, imagem_url, criado_em
FROM eventos
ORDER BY data;

-- name: BuscarEventoPorID :one
SELECT id, titulo, descricao, data, imagem_url, criado_em
FROM eventos
WHERE id = $1;

-- name: CriarEvento :one
INSERT INTO eventos (titulo, descricao, data, imagem_url)
VALUES ($1, $2, $3, $4)
RETURNING id, titulo, descricao, data, imagem_url, criado_em;

-- name: AtualizarEvento :one
UPDATE eventos
SET titulo = $2, descricao = $3, data = $4, imagem_url = $5
WHERE id = $1
RETURNING id, titulo, descricao, data, imagem_url, criado_em;

-- name: RemoverEvento :one
DELETE FROM eventos
WHERE id = $1
RETURNING id;

-- name: ListarFotosPorAluno :many
SELECT id, aluno_id, url, legenda, enviado_em
FROM fotos
WHERE aluno_id = $1
ORDER BY id;

-- name: CriarFoto :one
INSERT INTO fotos (aluno_id, url, legenda)
VALUES ($1, $2, $3)
RETURNING id, aluno_id, url, legenda, enviado_em;

-- name: BuscarAdminPorUsuario :one
SELECT id, usuario, senha_hash, role, criado_em
FROM admins
WHERE usuario = $1;

-- name: BuscarAdminPorID :one
SELECT id, usuario, senha_hash, role, criado_em
FROM admins
WHERE id = $1;

-- name: CriarAdmin :one
INSERT INTO admins (usuario, senha_hash, role)
VALUES ($1, $2, $3)
RETURNING id, usuario, senha_hash, role, criado_em;

-- name: CriarRefreshToken :exec
INSERT INTO refresh_tokens (admin_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: BuscarRefreshTokenAtivo :one
SELECT id, admin_id, token_hash, expires_at, revoked_at, created_at
FROM refresh_tokens
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- name: RevogarRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = $1
  AND revoked_at IS NULL;
