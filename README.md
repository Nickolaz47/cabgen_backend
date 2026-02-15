# Cabgen Backend

[English Version (Versão em Inglês)](./README.en.md)

Backend da plataforma **CABGen**, desenvolvido em **Go** utilizando o framework **Gin**.
Este projeto é uma reescrita do backend original do site [CABGen](https://cabgen.fiocruz.br/pt), com foco em desempenho, manutenibilidade e organização de código.

## Índice

1. [Tecnologias](#tecnologias)
2. [Estrutura do Projeto](#estrutura-do-projeto)
3. [Instalação](#instalação)
4. [Configuração](#configuração)
5. [Executando a API](#executando-a-api)
6. [Endpoints](#endpoints)
7. [Internacionalização (i18n)](#internacionalização-i18n)

## Tecnologias

- [Go](https://go.dev/) `>= 1.23.0`
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [go-i18n](https://github.com/nicksnyder/go-i18n)

## Estrutura do Projeto

```bash
.
├── cmd/                     # Ponto de entrada da aplicação
│   └── server/
│       └── main.go          # Inicialização da API
├── internal/                # Código interno (não exportável)
│   ├── auth/                # Autenticação (JWT e Cookies)
│   ├── config/              # Carregamento das variáveis de ambiente
│   ├── container/           # Inicialização de repositórios, services e handlers
│   ├── data/                # Dados estáticos (ex: countries.json)
│   ├── db/                  # Configuração e conexão com o banco
│   ├── email/               # Envio e configuração de emails
│   ├── events/              # Gerenciamento de eventos dentro da API
│   ├── handlers/            # Controllers (Gin)
│   ├── logging/             # Configuração e controle de logs
│   ├── middlewares/         # Middlewares da aplicação
│   ├── models/              # Models e mapeamento do banco
│   ├── repositories/        # Acesso e queries ao banco de dados
│   ├── responses/           # Padronização de respostas HTTP
│   ├── routes/              # Definição das rotas/endpoints
│   ├── security/            # Criptografia e hashing de senhas
│   ├── services/            # Regras de negócio
│   ├── testutils/           # Utilitários para testes
│   ├── translation/         # Internacionalização (i18n)
│   ├── utils/               # Funções utilitárias
│   └── validations/         # Validação de entradas
├── go.mod
├── go.sum
└── README.md
```

## Instalação

### Pré-requisitos

- [Go](https://go.dev/dl/) `>= 1.23.0`
- [PostgreSQL](https://www.postgresql.org/download/)
- [SQLite](https://sqlite.org/) (utilizado nos testes)

### Passos

```bash
git clone https://github.com/Nickolaz47/cabgen_backend.git
cd cabgen_backend
go mod tidy
```

## Configuração

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
# Banco de dados
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=

# JWT
SECRET_ACCESS_KEY=
SECRET_REFRESH_KEY=

# Frontend
FRONTEND_URL=          # Ex: http://localhost:3000

# API
PORT=                  # Ex: 8080
ENVIRONMENT=           # dev | prod
API_HOST=              # Ex: http://localhost:8080

# Usuário administrador padrão
ADMIN_PASSWORD=

# Configuração de email
SENDER_EMAIL=
SENDER_PASSWORD=
SMTP_HOST=
SMTP_PORT=
```

## Executando a API

### Ambiente de Desenvolvimento

O projeto utiliza **Air** para hot reload.

#### Instalação do Air

```bash
go install github.com/cosmtrek/air@latest
```

#### Execução

```bash
air
```

No arquivo `air.toml`, certifique-se de que o comando de build está configurado corretamente:

```toml
[build]
cmd = "go build -o ./tmp/main ./cmd/server/main.go"
```

### Ambiente de Produção

#### Execução Manual

1. Compile o binário:

```bash
go build -o cabgen-backend ./cmd/server
```

2. Execute a aplicação:

```bash
./cabgen-backend
```

#### Docker

1. Após configurar o `.env`, suba o compose:

```bash
docker compose up -d
```

## Internacionalização (i18n)

Idiomas suportados:

- pt-BR
- en-US
- es-ES

O idioma é detectado via header `Accept-Language`. Caso ele não seja enviado, o idioma padrão será o en-US.

### Comportamento em listagens e buscas

Para alguns recursos que possuem dados traduzidos (como **origens** e **fontes da amostra**), o idioma informado influencia diretamente o resultado das operações de **listagem** e **busca**.

Nesses casos:

- Apenas a tradução correspondente ao idioma solicitado será retornada;
- As demais traduções não são incluídas na resposta;
- As buscas textuais consideram exclusivamente o idioma ativo.

## Formato das Respostas e Códigos HTTP

A API utiliza um formato de resposta padronizado, composto pelos seguintes campos:

```json
{
  "data": {},
  "message": "",
  "error": ""
}
```

### Campos da Resposta

- **data**
  Utilizado para retornar dados da API.
  Está presente nos seguintes casos:
  - Respostas de leitura (`GET`)
  - Criação de recursos (`POST`)
  - Atualização de recursos (`PUT`)

- **message**
  Utilizado para mensagens informativas de sucesso.
  Está presente principalmente em:
  - Criação de recursos (`POST`)
  - Remoção de recursos (`DELETE`)

- **error**
  Presente **exclusivamente** quando ocorre algum erro durante o processamento da requisição.
  Contém uma mensagem descritiva do problema.

### Comportamento por Método HTTP

| Método | Campos retornados |
| ------ | ----------------- |
| GET    | `data`            |
| POST   | `data`, `message` |
| PUT    | `data`            |
| DELETE | `message`         |

### Códigos de Status HTTP

A API utiliza os seguintes códigos de status HTTP:

| Código | Descrição                                                                      |
| ------ | ------------------------------------------------------------------------------ |
| 200    | Requisição processada com sucesso                                              |
| 201    | Recurso criado com sucesso                                                     |
| 400    | Entrada inválida ou parâmetro de rota em formato incorreto (ex: UUID inválido) |
| 401    | Requisição sem token de autenticação                                           |
| 403    | Usuário desativado ou token de acesso expirado                                 |
| 404    | Recurso não encontrado                                                         |
| 409    | Tentativa de criação de recurso duplicado                                      |
| 500    | Erro interno inesperado                                                        |

## Endpoints

Os endpoints estão organizados em três níveis de acesso:

- **Público**: não requer autenticação
- **Common**: requer autenticação
- **Admin**: acesso restrito a administradores

### Público

#### Health Check

| Método | Endpoint      | Descrição                |
| ------ | ------------- | ------------------------ |
| GET    | `/api/health` | Verifica o status da API |

#### Autenticação

| Método | Endpoint             | Descrição                                 |
| ------ | -------------------- | ----------------------------------------- |
| POST   | `/api/auth/register` | Cadastro de usuário (necessita ativação)  |
| POST   | `/api/auth/login`    | Login e retorno de tokens JWT via cookies |
| POST   | `/api/auth/logout`   | Logout do usuário                         |
| POST   | `/api/auth/refresh`  | Renovação do token de acesso              |

#### Países

| Método | Endpoint               | Descrição                  |
| ------ | ---------------------- | -------------------------- |
| GET    | `/api/countries`       | Lista todos os países      |
| GET    | `/api/countries/:code` | Retorna um país específico |

### Common

#### Usuário

| Método | Endpoint        | Descrição                    |
| ------ | --------------- | ---------------------------- |
| GET    | `/api/users/me` | Dados do usuário autenticado |
| PUT    | `/api/users/me` | Atualiza dados do usuário    |

#### Origem

| Método | Endpoint       | Descrição            |
| ------ | -------------- | -------------------- |
| GET    | `/api/origins` | Lista origens ativas |

#### Sequenciador

| Método | Endpoint          | Descrição                   |
| ------ | ----------------- | --------------------------- |
| GET    | `/api/sequencers` | Lista sequenciadores ativos |

#### Fonte da Amostra

| Método | Endpoint              | Descrição                      |
| ------ | --------------------- | ------------------------------ |
| GET    | `/api/sample-sources` | Lista fontes de amostra ativas |

#### Laboratório

| Método | Endpoint            | Descrição                 |
| ------ | ------------------- | ------------------------- |
| GET    | `/api/laboratories` | Lista laboratórios ativos |

#### Microrganismo

| Método | Endpoint              | Descrição                   |
| ------ | --------------------  | -------------------------   |
| GET    | `/api/microorganisms` | Lista microrganismos ativos |

### Admin

Os endpoints administrativos seguem o padrão CRUD completo para **Usuários**, **Origens**, **Sequenciadores**, **Fontes da Amostra**, **Laboratórios** e **Microorganismos**:

#### Usuário

| Método | Endpoint                                | Descrição                     |
| ------ | --------------------------------------- | ----------------------------- |
| GET    | `/api/admin/users`                      | Lista todos os usuários       |
| GET    | `/api/admin/users/:id`                  | Retorna um usuário específico |
| POST   | `/api/admin/users`                      | Cria um usuário já ativado    |
| PUT    | `/api/admin/users/:id`                  | Atualiza um usuário           |
| PATCH  | `/api/admin/users/activate/:id`         | Ativa um usuário              |
| PATCH  | `/api/admin/users/deactivate/:id`       | Desativa um usuário           |
| DELETE | `/api/admin/users/:id`                  | Deleta um usuário             |

#### Origem

| Método | Endpoint                       | Descrição                     |
| ------ | ------------------------------ | ----------------------------- |
| GET    | `/api/admin/origins`           | Lista todas as origens        |
| GET    | `/api/admin/origins/:id`       | Retorna uma origem específica |
| GET    | `/api/admin/origins/search`    | Procura origens pelo nome     |
| POST   | `/api/admin/origins`           | Cria uma nova origem          |
| PUT    | `/api/admin/origins/:id`       | Atualiza uma origem           |
| DELETE | `/api/admin/origins/:id`       | Deleta uma origem             |

#### Sequenciador

| Método | Endpoint                       | Descrição                                   |
| ------ | ------------------------------ | ------------------------------------------- |
| GET    | `/api/admin/sequencers`        | Lista todos os sequenciadores               |
| GET    | `/api/admin/sequencers/:id`    | Retorna um sequenciador específico          |
| GET    | `/api/admin/sequencers/search` | Procura sequenciadores pela marca ou modelo |
| POST   | `/api/admin/sequencers`        | Cria um novo sequenciador                   |
| PUT    | `/api/admin/sequencers/:id`    | Atualiza um sequenciador                    |
| DELETE | `/api/admin/sequencers/:id`    | Deleta um sequenciador                      |

#### Fonte da Amostra

| Método | Endpoint                           | Descrição                                    |
| ------ | ---------------------------------- | -------------------------------------------- |
| GET    | `/api/admin/sample-sources`        | Lista todas as fontes da amostra             |
| GET    | `/api/admin/sample-sources/:id`    | Retorna uma fonte da amostra específica      |
| GET    | `/api/admin/sample-sources/search` | Procura fontes da amostra pelo nome ou grupo |
| POST   | `/api/admin/sample-sources`        | Cria uma nova fonte da amostra               |
| PUT    | `/api/admin/sample-sources/:id`    | Atualiza uma fonte da amostra                |
| DELETE | `/api/admin/sample-sources/:id`    | Deleta uma fonte da amostra                  |

#### Laboratório

| Método | Endpoint                         | Descrição                                    |
| ------ | -------------------------------- | -------------------------------------------- |
| GET    | `/api/admin/laboratories`        | Lista todos os laboratórios                  |
| GET    | `/api/admin/laboratories/:id`    | Retorna um laboratório específico            |
| GET    | `/api/admin/laboratories/search` | Procura laboratórios pelo nome ou abreviação |
| POST   | `/api/admin/laboratories`        | Cria um novo laboratório                     |
| PUT    | `/api/admin/laboratories/:id`    | Atualiza um laboratório                      |
| DELETE | `/api/admin/laboratories/:id`    | Deleta um laboratório                        |

#### Microrganismo

| Método | Endpoint                           | Descrição                                    |
| ------ | ---------------------------------- | -------------------------------------------- |
| GET    | `/api/admin/microorganisms`        | Lista todos os microrganismos                |
| GET    | `/api/admin/microorganisms/:id`    | Retorna um microrganismo específico          |
| GET    | `/api/admin/microorganisms/search` | Procura microrganismos pelo nome ou grupo    |
| POST   | `/api/admin/microorganisms`        | Cria um novo microrganismo                   |
| PUT    | `/api/admin/microorganisms/:id`    | Atualiza um microrganismo                    |
| DELETE | `/api/admin/microorganisms/:id`    | Deleta um microrganismo                      |

## TODO

- [x] Implementar logger nos services;
- [x] Modelar Microorganism;
- [ ] Modelar HealthService;
- [ ] Modelar Sample;
- [ ] Adicionar um volume para armazenar o events.db e as amostras recebidas;
