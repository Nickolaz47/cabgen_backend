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

- [Go](https://go.dev/) `>= 1.24.0`
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [go-i18n](https://github.com/nicksnyder/go-i18n)

## Estrutura do Projeto

### Estrutura das Pastas

```bash
.
├── cmd/                     # Ponto de entrada da aplicação
│   └── server/
│       └── main.go          # Inicialização da API
├── internal/                # Código interno (não exportável)
│   ├── auth/                # Autenticação (JWT e Cookies)
│   ├── config/              # Carregamento das variáveis de ambiente
│   ├── container/           # Inicialização de repositórios, services e handlers
│   ├── db/                  # Configuração e conexão com o banco
│   ├── email/               # Envio e configuração de emails
│   ├── handlers/            # Controllers (Gin)
│   ├── logging/             # Configuração e controle de logs
│   ├── middlewares/         # Middlewares da aplicação
│   ├── models/              # Models e mapeamento do banco
│   ├── pipeline/            # Processamento das amostras
│   ├── queue/               # Definição de tarefas e workers assíncronos (Redis/asynq)
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

### Estrutura do Código

O código segue a arquitetura em camadas, onde cada camada tem sua responsabilidade. A camada base é a camada de models, que é responsável por mapear os dados do banco de dados. A camada de repositories é responsável por acessar e consultar o banco de dados. A camada de services é responsável por implementar as regras de negócio. A camada de handlers é responsável por receber as requisições HTTP e retornar as respostas. Por sua vez, a camada de routes é responsável por definir as rotas/endpoints.

Logo para a construção de novos modelos, é necessário seguir essa ordem:
model -> repository -> service -> handler -> route

## Instalação

### Pré-requisitos

- [Go](https://go.dev/dl/) `>= 1.24.0`
- [PostgreSQL](https://www.postgresql.org/download/)
- [SQLite](https://sqlite.org/) (utilizado nos testes)

### Passos

```bash
git clone [https://github.com/Nickolaz47/cabgen_backend.git](https://github.com/Nickolaz47/cabgen_backend.git)
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
APP_ROOT=              # Sobrescreve o diretório raiz do projeto (auto-detectado se vazio)

# Usuário administrador padrão
ADMIN_PASSWORD=

# Configuração de email
SENDER_EMAIL=
SENDER_PASSWORD=
SMTP_HOST=
SMTP_PORT=

# Redis URL
REDIS_URL=

# Worker de Análise — Caminhos das ferramentas bioinformáticas (opcionais)
FASTQC_PATH=
UNICYCLER_PATH=
SPADES_PATH=
CHECKM_PATH=
KRAKEN2_PATH=
KRAKEN_DB_PATH=
FASTANI_PATH=
ABRICATE_PATH=
MLST_PATH=
RESFINDER_DB_PATH=

# Worker de Análise — Caminhos dos bancos de dados (opcionais)
POLI_DB_PSEUDO=
POLI_DB_KLEB=
POLI_DB_ENTERO=
POLI_DB_ACINETO=
OTHER_DB_PSEUDO=
OTHER_DB_KLEB=
OTHER_DB_ENTERO=
OTHER_DB_ACINETO=
FASTANI_LIST_KLEB=
FASTANI_LIST_ENTERO=
FASTANI_LIST_ACINETO=
```

## Executando a API

### Ambiente de Desenvolvimento

O projeto utiliza **Air** para hot reload.

#### Instalação do Air

```bash
go install [github.com/cosmtrek/air@latest](https://github.com/cosmtrek/air@latest)
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

#### Podman (rootless)

No podman rootless, o mapeamento de UIDs entre o container e o host é diferente do Docker. Para que `./logs` e `./uploads` fiquem acessíveis sem `sudo`, utilize o arquivo de override:

```bash
podman-compose -f docker-compose.yaml -f docker-compose.podman.yaml up -d
```

Para encurtar, adicione ao seu `~/.bashrc` ou `~/.zshrc`:

```bash
alias pdc='podman-compose -f docker-compose.yaml -f docker-compose.podman.yaml'
```

E então:

```bash
pdc up -d
```

O arquivo `docker-compose.podman.yaml` aplica `userns_mode: keep-id` apenas nos serviços da aplicação, preservando o comportamento padrão do postgres e redis.

## Seed

Na inicialização, a API executa o seed automático das seguintes entidades a partir dos arquivos JSON em `jsons/`:

- `countries.json`
- `microorganisms.json`
- `origins.json`
- `sequencers.json`
- `laboratories.json`
- `sample_sources.json`
- `health_services.json`

Cada tabela é populada apenas se estiver vazia. Os arquivos JSON podem ser mantidos fora do controle de versão via `.gitignore`.

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

#### **data**
Utilizado para retornar dados da API.
Está presente nos seguintes casos:

- Respostas de leitura (`GET`)
- Criação de recursos (`POST`)
- Atualização de recursos (`PUT`)

#### **message**

Utilizado para mensagens informativas de sucesso.
Está presente principalmente em:

- Criação de recursos (`POST`)
- Remoção de recursos (`DELETE`)

#### **error**

Presente **exclusivamente** quando ocorre algum erro durante o processamento da requisição.
Contém uma mensagem descritiva do problema.

### Comportamento por Método HTTP

| Método | Campos retornados |
| --- | --- |
| GET | `data` |
| POST | `data`, `message` |
| PUT | `data` |
| DELETE | `message` |

### Códigos de Status HTTP

A API utiliza os seguintes códigos de status HTTP:

| Código | Descrição |
| --- | --- |
| 200 | Requisição processada com sucesso |
| 201 | Recurso criado com sucesso |
| 400 | Entrada inválida ou parâmetro de rota em formato incorreto (ex: UUID inválido) |
| 401 | Requisição sem token de autenticação |
| 403 | Usuário desativado ou token de acesso expirado |
| 404 | Recurso não encontrado |
| 409 | Tentativa de criação de recurso duplicado |
| 410 | Recurso válido, mas não encontrado (apagado) |
| 500 | Erro interno inesperado |

## Endpoints

Os endpoints estão organizados em três níveis de acesso:

- **Público**: não requer autenticação
- **Common**: requer autenticação
- **Admin**: acesso restrito a administradores

### Público

#### Health Check

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/health` | Verifica o status da API |

#### Autenticação

| Método | Endpoint | Descrição |
| --- | --- | --- |
| POST | `/api/auth/register` | Cadastro de usuário (necessita ativação) |
| POST | `/api/auth/login` | Login e retorno de tokens JWT via cookies |
| POST | `/api/auth/refresh` | Renovação do token de acesso |
| POST | `/api/auth/forgot-password` | Solicitação de redefinição de senha |
| POST | `/api/auth/reset-password` | Redefinição de senha |

#### Países

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/countries` | Lista todos os países |
| GET | `/api/countries/:code` | Retorna um país específico |

#### Contato

| Método | Endpoint | Descrição |
| --- | --- | --- |
| POST | `/api/contact` | Cria um ticket de contato |

#### Métricas

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/metrics` | Retorna métricas gerais da plataforma (total de amostras, países, espécies e genes de resistência) |

### Common

#### Auth

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/auth/me` | Retorna os dados do usuário autenticado |
| POST | `/api/auth/logout` | Logout do usuário |

#### Usuário

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/users/me` | Dados do usuário autenticado |
| PUT | `/api/users/me` | Atualiza dados do usuário |
| POST | `/api/users/me/update-password` | Atualiza a senha do usuário autenticado |
| POST | `/api/users/me/request-email-update` | Solicita a atualização de e-mail e envia um link de confirmação |
| POST | `/api/users/me/confirm-email-update` | Confirma a atualização de e-mail usando o token do link |
| DELETE | `/api/users/me` | Deleta a conta do usuário autenticado |

#### Amostra

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/samples` | Lista todas as amostras do usuário |
| GET | `/api/samples/:sampleId` | Retorna uma amostra específica |
| POST | `/api/samples` | Cria uma nova amostra |
| PUT | `/api/samples/:sampleId/upload` | Faz upload dos arquivos (FASTQ/FASTA) |
| PUT | `/api/samples/:sampleId` | Atualiza os dados de uma amostra |
| DELETE | `/api/samples/:sampleId` | Deleta uma amostra |

#### Análise

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/analyses` | Lista todas as análises do usuário |
| GET | `/api/analyses/:analysisId` | Retorna uma análise específica |
| GET | `/api/analyses/:analysisId/download/zip` | Faz o download do arquivo ZIP da análise |
| POST | `/api/analyses` | Cria e inicia uma nova análise |
| POST | `/api/analyses/download/tsv` | Faz o download em lote (TSV) |
| DELETE | `/api/analyses/:analysisId` | Deleta uma análise |

#### Select Options

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/select-options/enum` | Retorna os enums para os selects do frontend (papéis, táxons, gêneros, tipos de serviço de saúde e tipos de análise) |
| GET | `/api/select-options/form` | Retorna as entidades ativas para os selects do frontend (laboratórios, sequenciadores, serviços de saúde, origens, microrganismos e fontes da amostra) |

#### Cidades

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/cities` | Retorna as cidades brasileiras para o select do frontend |

### Admin

Os endpoints administrativos seguem o padrão CRUD completo para **Usuários**, **Origens**, **Sequenciadores**, **Fontes da Amostra**, **Laboratórios**, **Microorganismos**, **Serviços de Saúde**, **Amostras**, **Análises** e **Tickets**:

#### Usuário

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/users` | Lista todos os usuários |
| GET | `/api/admin/users/:id` | Retorna um usuário específico |
| POST | `/api/admin/users` | Cria um usuário já ativado |
| PUT | `/api/admin/users/:id` | Atualiza um usuário |
| PATCH | `/api/admin/users/:id/activate` | Ativa um usuário |
| PATCH | `/api/admin/users/:id/deactivate` | Desativa um usuário |
| DELETE | `/api/admin/users/:id` | Deleta um usuário |

#### Origem

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/origins` | Lista todas as origens |
| GET | `/api/admin/origins/:id` | Retorna uma origem específica |
| GET | `/api/admin/origins/search` | Procura origens pelo nome |
| POST | `/api/admin/origins` | Cria uma nova origem |
| PUT | `/api/admin/origins/:id` | Atualiza uma origem |
| DELETE | `/api/admin/origins/:id` | Deleta uma origem |

#### Sequenciador

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/sequencers` | Lista todos os sequenciadores |
| GET | `/api/admin/sequencers/:id` | Retorna um sequenciador específico |
| GET | `/api/admin/sequencers/search` | Procura sequenciadores pela marca ou modelo |
| POST | `/api/admin/sequencers` | Cria um novo sequenciador |
| PUT | `/api/admin/sequencers/:id` | Atualiza um sequenciador |
| DELETE | `/api/admin/sequencers/:id` | Deleta um sequenciador |

#### Fonte da Amostra

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/sample-sources` | Lista todas as fontes da amostra |
| GET | `/api/admin/sample-sources/:id` | Retorna uma fonte da amostra específica |
| GET | `/api/admin/sample-sources/search` | Procura fontes da amostra pelo nome ou grupo |
| POST | `/api/admin/sample-sources` | Cria uma nova fonte da amostra |
| PUT | `/api/admin/sample-sources/:id` | Atualiza uma fonte da amostra |
| DELETE | `/api/admin/sample-sources/:id` | Deleta uma fonte da amostra |

#### Laboratório

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/laboratories` | Lista todos os laboratórios |
| GET | `/api/admin/laboratories/:id` | Retorna um laboratório específico |
| GET | `/api/admin/laboratories/search` | Procura laboratórios pelo nome ou abreviação |
| POST | `/api/admin/laboratories` | Cria um novo laboratório |
| PUT | `/api/admin/laboratories/:id` | Atualiza um laboratório |
| DELETE | `/api/admin/laboratories/:id` | Deleta um laboratório |

#### Microrganismo

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/microorganisms` | Lista todos os microrganismos |
| GET | `/api/admin/microorganisms/:id` | Retorna um microrganismo específico |
| GET | `/api/admin/microorganisms/search` | Procura microrganismos pelo nome ou grupo |
| POST | `/api/admin/microorganisms` | Cria um novo microrganismo |
| PUT | `/api/admin/microorganisms/:id` | Atualiza um microrganismo |
| DELETE | `/api/admin/microorganisms/:id` | Deleta um microrganismo |

#### Serviços de Saúde

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/health-services` | Lista todos os serviços de saúde |
| GET | `/api/admin/health-services/:id` | Retorna um serviços de saúde específico |
| GET | `/api/admin/health-services/search` | Procura serviços de saúde pelo nome ou grupo |
| POST | `/api/admin/health-services` | Cria um novo serviços de saúde |
| PUT | `/api/admin/health-services/:id` | Atualiza um serviços de saúde |
| DELETE | `/api/admin/health-services/:id` | Deleta um serviços de saúde |

#### Amostra

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/samples` | Lista todas as amostras |
| GET | `/api/admin/samples/:sampleId` | Retorna uma amostra específica |
| GET | `/api/admin/samples/genders` | Retorna os gêneros válidos para amostras |
| POST | `/api/admin/samples` | Cria uma nova amostra |
| PUT | `/api/admin/samples/:sampleId/upload` | Faz upload dos arquivos (FASTQ/FASTA) |
| PUT | `/api/admin/samples/:sampleId` | Atualiza os dados de uma amostra |
| DELETE | `/api/admin/samples/:sampleId` | Deleta uma amostra |

#### Análise

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/analyses` | Lista todas as análises |
| GET | `/api/admin/analyses/:analysisId` | Retorna uma análise específica |
| GET | `/api/admin/analyses/:analysisId/download/zip` | Faz o download do arquivo ZIP da análise |
| POST | `/api/admin/analyses` | Cria e inicia uma nova análise |
| POST | `/api/admin/analyses/download/tsv` | Faz o download em lote (TSV) |
| PUT | `/api/admin/analyses/:analysisId` | Atualiza o status/resultados da análise |
| DELETE | `/api/admin/analyses/:analysisId` | Deleta uma análise |

#### Ticket

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/tickets` | Lista todos os tickets |
| GET | `/api/admin/tickets/:ticketId` | Retorna um ticket específico |
| PUT | `/api/admin/tickets/:ticketId/assign` | Atribui um ticket a um administrador |
| PUT | `/api/admin/tickets/:ticketId/resolve` | Resolve um ticket |
| DELETE | `/api/admin/tickets/:ticketId` | Deleta um ticket |

#### Métricas

| Método | Endpoint | Descrição |
| --- | --- | --- |
| GET | `/api/admin/metrics` | Retorna métricas gerais da plataforma (amostras, países, espécies, genes de resistência, usuários, análises por status, países mais frequentes e espécies) |

## Organização do Diretório de Uploads

O diretório de uploads é organizado da seguinte forma:

```bash
uploads/
└── users/
    └── {user_id}/
        └── samples/
            └── {sample_id}/
                ├── reads.fastq
                └── analyses/
                    └── {analysis_id}/
                        ├── qc/                    
                        │   ├── fastqc.html
                        │   └── fastqc.zip
                        ├── assembly/
                        │   ├── contigs.fasta           
                        │   ├── assembly.gfa
                        │   ├── coverage.json           
                        │   ├── checkm_report.tsv        
                        │   ├── species_id.tsv           
                        │   └── annotation/               
                        ├── amr/                       
                        │   ├── resfinder.tsv
                        │   ├── virulence.tsv
                        │   ├── plasmids.tsv
                        │   ├── mlst.tsv
                        │   └── mutations.json
                        ├── report/
                        │   └── summary.json        
                        ├── logs/
                        │   └── pipeline.log
```

- **`qc/`**: controle de qualidade dos reads brutos (FastQC).
- **`assembly/`**: tudo que deriva da montagem — contigs (Unicycler), cobertura, qualidade da montagem (CheckM), identificação de espécie (Kraken2/FastANI) e anotação (Prokka).
- **`amr/`**: resultados de resistência, virulência, plasmídeos, MLST e mutações pontuais (ABRicate + ResFinder/VFDB/PlasmidFinder, `mlst`, BLASTx).
- **`report/`**: relatório final consolidado com os resultados clinicamente relevantes.
- **`logs/`**: logs de execução do pipeline.

## TODO

- [x] Implementar logger nos services;
- [x] Modelar Microorganism (Model + Repository + Service + Handler + Tests);
- [x] Modelar HealthService (Model + Repository + Service + Handler + Tests);
- [x] Modelar Sample (Model + Repository + Service + Handler + Tests);
- [x] Modelar Analysis (Model + Repository + Service + Handler + Tests)
- [x] Adicionar rota para recuperar a senha
- [x] Adicionar rota para redefinir a senha
- [x] Adicionar rotas para select
- [x] API pública -> Postgres -> Redis -> Pipeline -> API privada -> Postgres
- [x] Adicionar cidade como select no cadastro da amostra;
- [ ] Mostrar no resultado as versões de cada programa;
- [x] Permitir o download de vários resultados;
- [ ] Migrar os dados do MongoDB para o Postgresql;
- [x] Integrar com a pipeline;
- [x] Integrar com o frontend;
