# Entrega Final / Sprint 3

## Como executar a API

Variaveis de ambiente:

| Variavel | Descricao | Exemplo |
|----------|-----------|---------|
| `DATABASE_URL` | Conexao PostgreSQL. Se ausente, a API usa memoria apenas para desenvolvimento local. | `postgres://anuario:anuario@localhost:5432/anuario?sslmode=disable` |
| `JWT_SECRET` | Segredo HMAC para assinar JWTs. Deve ter pelo menos 16 caracteres. | `uma-chave-local-segura` |
| `ADMIN_USER` | Usuario admin criado no banco ao iniciar. | `admin` |
| `ADMIN_PASSWORD` | Senha admin criada no banco ao iniciar. | `admin123` |
| `PORT` | Porta HTTP. | `8080` |

Comandos principais:

```bash
go mod download
go test ./... -v -count=1
go run ./cmd/api
```

Ao usar PostgreSQL, o schema em `interno/db/schema.sql` e aplicado automaticamente na inicializacao, e o admin padrao e criado se ainda nao existir.

## Endpoints principais

Frontend:

| Rota | Descricao |
|------|-----------|
| `/` | Pagina principal do anuario com alunos, galeria, eventos e painel admin |
| `/login.html` | Tela de login que consome `/api/v1/auth/login` |

API:

| Metodo | Rota | Descricao | Auth |
|--------|------|-----------|------|
| `POST` | `/api/v1/auth/login` | Login do admin, retorna JWT e refresh token | Publica |
| `POST` | `/api/v1/auth/refresh` | Rotaciona refresh token e retorna novo JWT | Publica com refresh token |
| `GET` | `/api/v1/alunos` | Lista alunos | Publica |
| `GET` | `/api/v1/alunos/{id}` | Detalhe do aluno com fotos aninhadas `1:N` | Publica |
| `GET` | `/api/v1/alunos/{id}/fotos` | Lista fotos de um aluno | Publica |
| `POST` | `/api/v1/alunos` | Cria aluno | Bearer JWT admin |
| `PUT` | `/api/v1/alunos/{id}` | Atualiza aluno | Bearer JWT admin |
| `DELETE` | `/api/v1/alunos/{id}` | Remove aluno | Bearer JWT admin |
| `POST` | `/api/v1/alunos/{id}/fotos` | Cria foto vinculada a aluno | Bearer JWT admin |
| `GET` | `/api/v1/eventos` | Lista eventos | Publica |
| `GET` | `/api/v1/eventos/{id}` | Detalhe do evento | Publica |
| `POST` | `/api/v1/eventos` | Cria evento | Bearer JWT admin |
| `PUT` | `/api/v1/eventos/{id}` | Atualiza evento | Bearer JWT admin |
| `DELETE` | `/api/v1/eventos/{id}` | Remove evento | Bearer JWT admin |

## Persistencia PostgreSQL + sqlc

- `sqlc.yaml` documenta a configuracao sqlc.
- `interno/db/schema.sql` define `admins`, `alunos`, `eventos`, `fotos` e `refresh_tokens`.
- `interno/db/query.sql` concentra as queries tipadas.
- `interno/db/*.sql.go` contem o codigo Go gerado/compatibilizado com sqlc.
- `interno/repository/postgres.go` usa o pacote `db` para persistir CRUD e o relacionamento `alunos 1:N fotos`.

## Seguranca implementada

- JWT HS256 no login e middleware `Bearer` para rotas protegidas.
- Autorizacao por papel `admin` nas rotas de escrita.
- Refresh token opaco com hash SHA-256 armazenado e rotacao no endpoint `/auth/refresh`.
- Rate limiting nos endpoints de autentificacao.
- Headers de seguranca: `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy` e `Content-Security-Policy`.
- Validacao e sanitizacao de entrada para alunos, eventos e fotos.
- Limite de tamanho do corpo JSON e rejeicao de campos desconhecidos.

## Testes e CI

- A suite local executa mais de 10 testes automatizados cobrindo controllers, auth, middlewares e repositorio em memoria.
- O workflow `.github/workflows/testes.yml` sobe PostgreSQL 16 e define `TEST_DATABASE_URL`, permitindo tambem teste de integracao do repositorio Postgres.
- Comando validado localmente: `go test ./... -v -count=1`.
