# Cabgen Backend

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

* [Go](https://go.dev/) `>= 1.23.0`
* [Gin](https://gin-gonic.com/)
* [GORM](https://gorm.io/)
* [PostgreSQL](https://www.postgresql.org/)
* [go-i18n](https://github.com/nicksnyder/go-i18n)

## Estrutura do Projeto

```bash
.
├── cmd/                     # Ponto de entrada da aplicação
│   └── server/
│       └── main.go          # Inicialização da API
├── internal/                # Código interno (não exportável)
│   ├── auth/                # Autenticação (JWT e Cookies)
│   ├── config/              # Carregamento das variáveis de ambiente
│   ├── container/           # Inicialização de services e handlers
│   ├── data/                # Dados estáticos (ex: countries.json)
│   ├── db/                  # Configuração e conexão com o banco
│   ├── email/               # Envio e configuração de emails
│   ├── handlers/            # Controllers (Gin)
│   ├── logging/             # Configuração e controle de logs
│   ├── middlewares/         # Middlewares da aplicação
│   ├── models/              # Models e mapeamento do banco
│   ├── repository/          # Acesso e queries ao banco de dados
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

* [Go](https://go.dev/dl/) `>= 1.23.0`
* [PostgreSQL](https://www.postgresql.org/download/)
* [SQLite](https://sqlite.org/) (utilizado nos testes)

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

1. Compile o binário:

```bash
go build -o cabgen-backend ./cmd/server
```

2. Execute a aplicação:

```bash
./cabgen-backend
```

## Internacionalização (i18n)

Idiomas suportados:

* pt-BR
* en-US
* es-ES

O idioma é detectado via header `Accept-Language`. Caso ele não seja enviado, o idioma padrão será o en-US.

### Comportamento em listagens e buscas

Para alguns recursos que possuem dados traduzidos (como **origens** e **fontes da amostra**), o idioma informado influencia diretamente o resultado das operações de **listagem** e **busca**.

Nesses casos:

* Apenas a tradução correspondente ao idioma solicitado será retornada;
* As demais traduções não são incluídas na resposta;
* As buscas textuais consideram exclusivamente o idioma ativo.

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

* **data**
  Utilizado para retornar dados da API.
  Está presente nos seguintes casos:

  * Respostas de leitura (`GET`)
  * Criação de recursos (`POST`)
  * Atualização de recursos (`PUT`)

* **message**
  Utilizado para mensagens informativas de sucesso.
  Está presente principalmente em:

  * Criação de recursos (`POST`)
  * Remoção de recursos (`DELETE`)

* **error**
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

* **Público**: não requer autenticação
* **Common**: requer autenticação
* **Admin**: acesso restrito a administradores

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

| Método | Endpoint             | Descrição                  |
| ------ | -------------------- | -------------------------- |
| GET    | `/api/country`       | Lista todos os países      |
| GET    | `/api/country/:code` | Retorna um país específico |

### Common

#### Usuário

| Método | Endpoint       | Descrição                    |
| ------ | -------------- | ---------------------------- |
| GET    | `/api/user/me` | Dados do usuário autenticado |
| PUT    | `/api/user/me` | Atualiza dados do usuário    |

#### Origem

| Método | Endpoint      | Descrição            |
| ------ | ------------- | -------------------- |
| GET    | `/api/origin` | Lista origens ativas |

#### Sequenciador

| Método | Endpoint         | Descrição                   |
| ------ | ---------------- | --------------------------- |
| GET    | `/api/sequencer` | Lista sequenciadores ativos |

#### Fonte da Amostra

| Método | Endpoint            | Descrição                      |
| ------ | ------------------- | ------------------------------ |
| GET    | `/api/sampleSource` | Lista fontes de amostra ativas |

#### Laboratório

| Método | Endpoint          | Descrição                 |
| ------ | ----------------- | ------------------------- |
| GET    | `/api/laboratory` | Lista laboratórios ativos |

### Admin

Os endpoints administrativos seguem o padrão CRUD completo para **Usuários**, **Origens**, **Sequenciadores**, **Fontes da Amostra** e **Laboratórios**:

#### Usuário

| Método | Endpoint                               | Descrição                     |
| ------ | -------------------------------------- | ----------------------------- |
| GET    | `/api/admin/user`                      | Lista todos os usuários       |
| GET    | `/api/admin/user/:username`            | Retorna um usuário específico |
| POST   | `/api/admin/user`                      | Cria um usuário já ativado    |
| PUT    | `/api/admin/user/:username`            | Atualiza um usuário           |
| PUT    | `/api/admin/user/activation/:username` | Ativa/desativa um usuário     |
| DELETE | `/api/admin/user/:username`            | Deleta um usuário             |

#### Origem

| Método | Endpoint                         | Descrição                     |
| ------ | -------------------------------- | ----------------------------- |
| GET    | `/api/admin/origin`              | Lista todas as origens        |
| GET    | `/api/admin/origin/:originId`    | Retorna uma origem específica |
| PUT    | `/api/admin/origin/search?name=` | Procura origens pelo nome     |
| POST   | `/api/admin/origin`              | Cria uma nova origem          |
| PUT    | `/api/admin/origin/:originId`    | Atualiza uma origem           |
| DELETE | `/api/admin/origin/:originId`    | Deleta uma origem             |

#### Sequenciador

| Método | Endpoint                                    | Descrição                                   |
| ------ | ------------------------------------------- | ------------------------------------------- |
| GET    | `/api/admin/sequencer`                      | Lista todos os sequenciadores               |
| GET    | `/api/admin/sequencer/:sequencerId`         | Retorna um sequenciador específico          |
| PUT    | `/api/admin/sequencer/search?brandOrModel=` | Procura sequenciadores pela marca ou modelo |
| POST   | `/api/admin/sequencer`                      | Cria um novo sequenciador                   |
| PUT    | `/api/admin/sequencer/:sequencerId`         | Atualiza um sequenciador                    |
| DELETE | `/api/admin/sequencer/:sequencerId`         | Deleta um sequenciador                      |

#### Fonte da Amostra

| Método | Endpoint                                      | Descrição                                    |
| ------ | --------------------------------------------- | -------------------------------------------- |
| GET    | `/api/admin/sampleSource`                     | Lista todas as fontes da amostra             |
| GET    | `/api/admin/sampleSource/:sampleSourceId`     | Retorna uma fonte da amostra específica      |
| PUT    | `/api/admin/sampleSource/search?nameOrGroup=` | Procura fontes da amostra pelo nome ou grupo |
| POST   | `/api/admin/sampleSource`                     | Cria uma nova fonte da amostra               |
| PUT    | `/api/admin/sampleSource/:sampleSourceId`     | Atualiza uma fonte da amostra                |
| DELETE | `/api/admin/sampleSource/:sampleSourceId`     | Deleta uma fonte da amostra                  |

#### Laboratório

| Método | Endpoint                                           | Descrição                                    |
| ------ | -------------------------------------------------- | -------------------------------------------- |
| GET    | `/api/admin/laboratory`                            | Lista todos os laboratórios                  |
| GET    | `/api/admin/laboratory/:laboratoryId`              | Retorna um laboratório específico            |
| PUT    | `/api/admin/laboratory/search?nameOrAbbreviation=` | Procura laboratórios pelo nome ou abreviação |
| POST   | `/api/admin/laboratory`                            | Cria um novo laboratório                     |
| PUT    | `/api/admin/laboratory/:laboratoryId`              | Atualiza um laboratório                      |
| DELETE | `/api/admin/laboratory/:laboratoryId`              | Deleta um laboratório                        |
