# Cabgen Backend

API desenvolvida utilizando a linguagem [Go](https://go.dev/) juntamente com o framework [Gin](https://gin-gonic.com/en/docs/). Esta API é uma nova versão do backend do site do [CABGen](https://cabgen.fiocruz.br/pt).

## Índice

1. [Tecnologias](#tecnologias)
2. [Estrutura do projeto](#estrutura-do-projeto)
3. [Instalação](#instalação)
4. [Configuração](#configuração)
5. [Executando a API](#executando-a-api)
6. [Endpoints](#endpoints)
7. [Internacionalização (i18n)](#internacionalização-i18n)

## Tecnologias

- [Go](https://golang.org/) 1.23.0
- [Gin](https://gin-gonic.com/)
- [go-i18n](https://github.com/nicksnyder/go-i18n)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)

## Estrutura do Projeto

```bash
.
├── cmd/                     # Ponto de entrada da aplicação
│   └── server/
│       └── main.go           # Inicializa a API
├── internal/                 # Código interno (não exportável)
│   ├── auth/                 # JWT e cookies
│   ├── config/               # Carregamento das variáveis de ambiente
│   ├── data/                 # Dados estáticos (ex: countries.json)
│   ├── db/                   # Configuração e conexão com o banco de dados
│   ├── handlers/             # Controladores Gin
│   ├── logging/              # Controle de logs
│   ├── middlewares/          # Middlewares
│   ├── models/               # Estruturas de dados e mapeamento do banco
│   ├── repository            # Queries do banco de dados
│   ├── responses/            # Padronização de respostas
│   ├── routes/               # Definição de endpoints
│   ├── security/             # Criptografia de senhas
│   ├── testutils/            # Funções auxiliares para os testes
│   ├── translation/          # Arquivos e lógica de i18n
│   ├── utils/                # Funções auxiliares
│   └── validations/          # Validação de entrada e regras de negócio
├── go.mod
├── go.sum
└── README.md
```

## Instalação

1. Instale o [Go](https://go.dev/dl/) (versão 1.23.0 ou superior recomendada).
2. Instale o [PostgreSQL](https://www.postgresql.org/download/) e configure seu banco de dados.
3. Instale o [SQLite](https://sqlite.org/) para os testes.

Em seguida, clone este repositório e baixe as dependências:

```bash
git clone https://github.com/Nickolaz47/cabgen_backend.git
cd cabgen_backend
go mod tidy
```

## Configuração

```env
# Banco de dados
DB_USER=              # Usuário do banco PostgreSQL
DB_PASSWORD=          # Senha do banco
DB_NAME=              # Nome do banco

# Token JWT
SECRET_ACCESS_KEY=    # Chave secreta para assinar tokens de acesso
SECRET_REFRESH_KEY=   # Chave secreta para assinar tokens de refresh

# Frontend
FRONTEND_URL=         # URL do frontend (ex: http://localhost:3000)

# API
PORT=                 # Porta da API (ex: 8080)
ENVIRONMENT=          # Ambiente de execução: dev | prod
API_HOST=             # URL base da API (ex: http://localhost:8080)

# Usuário administrador padrão
ADMIN_PASSWORD=       # Senha inicial do admin
```

## Executando a API

### Desenvolvimento

Durante o desenvolvimento, o Air é utilizado para hot reload, facilitando o processo de testar alterações sem reiniciar manualmente a aplicação.

Para rodar a API em modo desenvolvimento, execute:

```bash
air
```

Certifique-se de que o air está instalado globalmente. Caso não esteja, instale com:

```bash
go install github.com/cosmtrek/air@latest
```

Altere o caminho do build no cmd dentro do arquivo `air.toml`:

```toml
[build]
cmd = "go build -o ./tmp/main ./cmd/server/main.go"
```

### Produção

Para executar a API em produção, siga os passos:

1. Compile o binário:

```bash
go build -o cabgen-backend ./cmd/server
```

2. Execute o binário gerado:

```bash
./cabgen-backend
```

Por padrão, a aplicação irá usar as configurações do arquivo .env e escutará na porta configurada (PORT).

## Endpoints

### Público

| Método | Endpoint           | Descrição                       |
|--------|--------------------|--------------------------------|
| GET   | `/api/health`| Verifica se a API está online    |

### Autenticação

| Método | Endpoint           | Descrição                       |
|--------|--------------------|--------------------------------|
| POST   | `/api/auth/register`| Cadastra o usuário que precisa ser ativado por um admin |
| POST   | `/api/auth/login` | Faz login e retorna tokens JWT via Cookies |
| POST   | `/api/auth/logout`| Encerra a sessão do usuário     |
| POST   | `/api/auth/refresh`| Renova o token de acesso    |

### Países

| Método | Endpoint               | Descrição                          |
|--------|------------------------|-----------------------------------|
| GET    | `/api/country`    | Retorna todos os países   |
| GET    | `/api/country/:code`    | Retorna um país específico        |

### Usuários

| Método | Endpoint               | Descrição                          |
|--------|------------------------|-----------------------------------|
| GET    | `/api/user/me`    | Retorna dados do usuário logado   |
| PUT    | `/api/user/me`    | Atualiza dados do usuário logado         |

### Admin

| Método | Endpoint                 | Descrição                              |
|--------|--------------------------|---------------------------------------|
| GET    | `/api/admin/user`       | Lista todos os usuários                |
| GET    | `/api/admin/user/:username`       | Retorna um usuário específico       |
| POST | `/api/admin/user`   | Cria um usuário já ativado                 |
| PUT | `/api/admin/user/:username`   | Atualiza um usuário                |
| PUT | `/api/admin/user/activation/:username`   | Ativa/desativa um usuário                |
| DELETE | `/api/admin/user/:username`   | Deleta um usuário                |

## Internacionalização (i18n)

Idiomas suportados:

- pt-BR
- en-US
- es-ES

O idioma é detectado via header `Accept-Language`. Caso ele não seja enviado, o idioma padrão será o en-US.
