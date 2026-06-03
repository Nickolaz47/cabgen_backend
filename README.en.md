Here is the updated English version of your README, synchronized with the latest changes made to the Portuguese version.

---

# Cabgen Backend

[Portuguese Version (Versão em Português)](./README.md)

Backend of the **CABGen** platform, developed in **Go** using the **Gin** framework.
This project is a rewrite of the original backend for the [CABGen](https://cabgen.fiocruz.br/pt) website, focusing on performance, maintainability, and code organization.

## Table of Contents

1. [Technologies](#technologies)
2. [Project Structure](#project-structure)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Running the API](#running-the-api)
6. [Endpoints](#endpoints)
7. [Internationalization (i18n)](#internationalization-i18n)

## Technologies

- [Go](https://go.dev/) `>= 1.23.0`
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [go-i18n](https://github.com/nicksnyder/go-i18n)

## Project Structure

### Folder Structure

```bash
.
├── cmd/                     # Application entry point
│   └── server/
│       └── main.go          # API initialization
├── internal/                # Internal code (non-exportable)
│   ├── auth/                # Authentication (JWT and Cookies)
│   ├── config/              # Environment variables loading
│   ├── container/           # Repositories, services, and handlers initialization
│   ├── data/                # Static data (e.g., countries.json)
│   ├── db/                  # Database configuration and connection
│   ├── email/               # Email sending and configuration
│   ├── events/              # Event management within the API
│   ├── handlers/            # Controllers (Gin)
│   ├── logging/             # Logging configuration and control
│   ├── middlewares/         # Application middlewares
│   ├── models/              # Models and database mapping
│   ├── repositories/        # Database access and queries
│   ├── responses/           # HTTP response standardization
│   ├── routes/              # Route and endpoint definition
│   ├── security/            # Password encryption and hashing
│   ├── services/            # Business logic
│   ├── testutils/           # Testing utilities
│   ├── translation/         # Internationalization (i18n)
│   ├── utils/               # Utility functions
│   └── validations/         # Input validation
├── go.mod
├── go.sum
└── README.md
```

### Code Structure

The code follows a layered architecture where each layer has its own responsibility. The base layer is **Models**, responsible for mapping database data. The **Repositories** layer handles database access and queries. The **Services** layer implements business logic. The **Handlers** layer receives HTTP requests and returns responses. Finally, the **Routes** layer defines the endpoints.

To build new features, the following order should be followed:
model -> repository -> service -> handler -> route

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) `>= 1.23.0`
- [PostgreSQL](https://www.postgresql.org/download/)
- [SQLite](https://sqlite.org/) (used in tests)

### Steps

```bash
git clone https://github.com/Nickolaz47/cabgen_backend.git
cd cabgen_backend
go mod tidy
```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Database
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=

# JWT
SECRET_ACCESS_KEY=
SECRET_REFRESH_KEY=

# Frontend
FRONTEND_URL=          # e.g., http://localhost:3000

# API
PORT=                  # e.g., 8080
ENVIRONMENT=           # dev | prod
API_HOST=              # e.g., http://localhost:8080

# Default administrator
ADMIN_PASSWORD=

# Email configuration
SENDER_EMAIL=
SENDER_PASSWORD=
SMTP_HOST=
SMTP_PORT=
```

## Running the API

### Development Environment

The project uses **Air** for hot reloading.

#### Air Installation

```bash
go install github.com/cosmtrek/air@latest
```

#### Execution

```bash
air
```

Ensure the build command in `air.toml` is correctly configured:

```toml
[build]
cmd = "go build -o ./tmp/main ./cmd/server/main.go"
```

### Production Environment

#### Manual Execution

1. Compile the binary:

```bash
go build -o cabgen-backend ./cmd/server
```

2. Run the application:

```bash
./cabgen-backend
```

#### Docker

1. After configuring the `.env`, start the containers:

```bash
docker compose up -d
```

## Internationalization (i18n)

Supported languages:

- pt-BR
- en-US
- es-ES

The language is detected via the `Accept-Language` header. If missing, the default is `en-US`.

### Behavior in Lists and Searches

For resources containing translated data (such as **origins** and **sample sources**), the requested language directly influences **list** and **search** results.

In these cases:

- Only the translation corresponding to the requested language is returned;
- Other translations are excluded from the response;
- Text searches consider only the active language.

## Response Format and HTTP Status Codes

The API uses a standardized response format:

```json
{
  "data": {},
  "message": "",
  "error": ""
}
```

### Response Fields

- **data**: Used to return API data. Present in `GET` responses, resource creation (`POST`), and updates (`PUT`).
- **message**: Used for informative success messages. Primarily present in resource creation (`POST`) and deletion (`DELETE`).
- **error**: Present **exclusively** when an error occurs. Contains a descriptive message of the problem.

### Behavior by HTTP Method

| Method | Fields Returned   |
| ------ | ----------------- |
| GET    | `data`            |
| POST   | `data`, `message` |
| PUT    | `data`            |
| DELETE | `message`         |

### HTTP Status Codes

| Code | Description                                                               |
| ---- | ------------------------------------------------------------------------- |
| 200  | Request processed successfully                                            |
| 201  | Resource created successfully                                             |
| 400  | Invalid input or route parameter in wrong format (e.g., invalid UUID)     |
| 401  | Request missing authentication token                                      |
| 403  | User disabled or access token expired                                     |
| 404  | Resource not found                                                        |
| 409  | Attempt to create a duplicate resource                                    |
| 500  | Unexpected internal error                                                 |

## Endpoints

Endpoints are organized into three access levels:

- **Public**: No authentication required.
- **Common**: Requires authentication.
- **Admin**: Restricted to administrators.

### Public

#### Health Check

| Method | Endpoint      | Description           |
| ------ | ------------- | --------------------- |
| GET    | `/api/health` | Checks the API status |

#### Authentication

| Method | Endpoint             | Description                                 |
| ------ | -------------------- | ------------------------------------------- |
| POST   | `/api/auth/register` | User registration (requires activation)     |
| POST   | `/api/auth/login`    | Login and returns JWT tokens via cookies    |
| POST   | `/api/auth/logout`   | User logout                                 |
| POST   | `/api/auth/refresh`  | Access token renewal                        |

#### Countries

| Method | Endpoint               | Description                |
| ------ | ---------------------- | ------------------------   |
| GET    | `/api/countries`       | Lists all countries        |
| GET    | `/api/countries/:code` | Returns a specific country |

### Common

#### User

| Method | Endpoint        | Description                     |
| ------ | --------------- | -----------------------------   |
| GET    | `/api/users/me` | Authenticated user data         |
| PUT    | `/api/users/me` | Updates authenticated user data |

#### Origin

| Method | Endpoint       | Description         |
| ------ | -------------- | ------------------- |
| GET    | `/api/origins` | Lists active origins|

#### Sequencer

| Method | Endpoint          | Description            |
| ------ | ----------------- | ---------------------- |
| GET    | `/api/sequencers` | Lists active sequencers|

#### Sample Source

| Method | Endpoint              | Description                |
| ------ | --------------------- | -------------------------- |
| GET    | `/api/sample-sources` | Lists active sample sources|

#### Laboratory

| Method | Endpoint            | Description               |
| ------ | ------------------- | ------------------------  |
| GET    | `/api/laboratories` | Lists active laboratories |

#### Microorganism

| Method | Endpoint              | Description                 |
| ------ | --------------------- | --------------------------  |
| GET    | `/api/microorganisms` | Lists active microorganisms |

#### Health Service

| Method | Endpoint               | Description                  |
| ------ | ---------------------- | --------------------------   |
| GET    | `/api/health-services` | Lists active health services |

### Admin

Administrative endpoints follow the full CRUD pattern for **Users**, **Origins**, **Sequencers**, **Sample Sources**, **Laboratories**, **Microorganisms**, and **Health Services**:

#### User

| Method | Endpoint                          | Description              |
| ------ | --------------------------------- | ------------------------ |
| GET    | `/api/admin/users`                | Lists all users          |
| GET    | `/api/admin/users/:id`            | Returns a specific user  |
| POST   | `/api/admin/users`                | Creates a pre-activated user |
| PUT    | `/api/admin/users/:id`            | Updates a user           |
| PATCH  | `/api/admin/users/activate/:id`   | Activates a user         |
| PATCH  | `/api/admin/users/deactivate/:id` | Deactivates a user       |
| DELETE | `/api/admin/users/:id`            | Deletes a user           |

#### Origin

| Method | Endpoint                    | Description              |
| ------ | --------------------------- | ------------------------ |
| GET    | `/api/admin/origins`        | Lists all origins        |
| GET    | `/api/admin/origins/:id`    | Returns a specific origin |
| GET    | `/api/admin/origins/search` | Searches origins by name |
| POST   | `/api/admin/origins`        | Creates a new origin     |
| PUT    | `/api/admin/origins/:id`    | Updates an origin        |
| DELETE | `/api/admin/origins/:id`    | Deletes an origin        |

#### Sequencer

| Method | Endpoint                       | Description                          |
| ------ | ------------------------------ | ------------------------------------ |
| GET    | `/api/admin/sequencers`        | Lists all sequencers                 |
| GET    | `/api/admin/sequencers/:id`    | Returns a specific sequencer          |
| GET    | `/api/admin/sequencers/search` | Searches sequencers by brand or model |
| POST   | `/api/admin/sequencers`        | Creates a new sequencer              |
| PUT    | `/api/admin/sequencers/:id`    | Updates a sequencer                  |
| DELETE | `/api/admin/sequencers/:id`    | Deletes a sequencer                  |

#### Sample Source

| Method | Endpoint                           | Description                            |
| ------ | ---------------------------------- | -------------------------------------- |
| GET    | `/api/admin/sample-sources`        | Lists all sample sources               |
| GET    | `/api/admin/sample-sources/:id`    | Returns a specific sample source        |
| GET    | `/api/admin/sample-sources/search` | Searches sample sources by name or group|
| POST   | `/api/admin/sample-sources`        | Creates a new sample source            |
| PUT    | `/api/admin/sample-sources/:id`    | Updates a sample source                |
| DELETE | `/api/admin/sample-sources/:id`    | Deletes a sample source                |

#### Laboratory

| Method | Endpoint                         | Description                                 |
| ------ | -------------------------------- | ------------------------------------------- |
| GET    | `/api/admin/laboratories`        | Lists all laboratories                      |
| GET    | `/api/admin/laboratories/:id`    | Returns a specific laboratory               |
| GET    | `/api/admin/laboratories/search` | Searches laboratories by name or abbreviation|
| POST   | `/api/admin/laboratories`        | Creates a new laboratory                    |
| PUT    | `/api/admin/laboratories/:id`    | Updates a laboratory                        |
| DELETE | `/api/admin/laboratories/:id`    | Deletes a laboratory                        |

#### Microorganism

| Method | Endpoint                           | Description                                 |
| ------ | ---------------------------------- | ------------------------------------------- |
| GET    | `/api/admin/microorganisms`        | Lists all microorganisms                    |
| GET    | `/api/admin/microorganisms/:id`    | Returns a specific microorganism            |
| GET    | `/api/admin/microorganisms/search` | Searches microorganisms by name or group    |
| POST   | `/api/admin/microorganisms`        | Creates a new microorganism                 |
| PUT    | `/api/admin/microorganisms/:id`    | Updates a microorganism                     |
| DELETE | `/api/admin/microorganisms/:id`    | Deletes a microorganism                     |

#### Health Service

| Method | Endpoint                            | Description                                 |
| ------ | ----------------------------------- | ------------------------------------------- |
| GET    | `/api/admin/health-services`        | Lists all health services                   |
| GET    | `/api/admin/health-services/:id`    | Returns a specific health service           |
| GET    | `/api/admin/health-services/search` | Searches health services by name or group   |
| POST   | `/api/admin/health-services`        | Creates a new health service                |
| PUT    | `/api/admin/health-services/:id`    | Updates a health service                    |
| DELETE | `/api/admin/health-services/:id`    | Deletes a health service                    |
