# Cabgen Backend

[Portuguese Version (Versão em Português)](./README.md)

Backend of the **CABGen** platform, developed in **Go** using the **Gin** framework.
This project is a rewrite of the original backend for the [CABGen](https://cabgen.fiocruz.br/pt) website, with focus on performance, maintainability, and code organization.

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

```bash
.
├── cmd/                     # Application entry point
│   └── server/
│       └── main.go          # API initialization
├── internal/                # Internal code (non-exportable)
│   ├── auth/                # Authentication (JWT and Cookies)
│   ├── config/              # Environment variables loading
│   ├── container/           # Services and handlers initialization
│   ├── data/                # Static data (ex: countries.json)
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
└── README.en.md
```

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
FRONTEND_URL=          # Ex: http://localhost:3000

# API
PORT=                  # Ex: 8080
ENVIRONMENT=           # dev | prod
API_HOST=              # Ex: http://localhost:8080

# Default admin user
ADMIN_PASSWORD=

# Email configuration
SENDER_EMAIL=
SENDER_PASSWORD=
SMTP_HOST=
SMTP_PORT=
```

## Running the API

### Development Environment

The project uses **Air** for hot reload.

#### Air Installation

```bash
go install github.com/cosmtrek/air@latest
```

#### Execution

```bash
air
```

In the `air.toml` file, ensure the build command is configured correctly:

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

1. After configuring the `.env`, start the compose:

```bash
docker compose up -d
```

## Internationalization (i18n)

Supported languages:

- pt-BR
- en-US
- es-ES

Language is detected via the `Accept-Language` header. If not provided, the default language is en-US.

### Behavior in lists and searches

For resources with translated data (such as **origins** and **sample sources**), the specified language directly influences **list** and **search** operations.

In these cases:

- Only the translation for the requested language is returned
- Other translations are not included in the response
- Text searches consider only the active language

## Response Format and HTTP Status Codes

The API uses a standardized response format with the following fields:

```json
{
  "data": {},
  "message": "",
  "error": ""
}
```

### Response Fields

- **data**: Used to return API data. Present in: GET requests, resource creation (POST), resource updates (PUT)
- **message**: Used for success informative messages. Mainly present in: resource creation (POST), resource deletion (DELETE)
- **error**: Present **exclusively** when an error occurs during request processing. Contains a descriptive problem message.

### Behavior by HTTP Method

| Method | Fields returned   |
| ------ | ----------------- |
| GET    | `data`            |
| POST   | `data`, `message` |
| PUT    | `data`            |
| DELETE | `message`         |

### HTTP Status Codes

| Code | Description                                                           |
| ---- | --------------------------------------------------------------------- |
| 200  | Request processed successfully                                        |
| 201  | Resource created successfully                                         |
| 400  | Invalid input or route parameter in wrong format (e.g., invalid UUID) |
| 401  | Request without authentication token                                  |
| 403  | User disabled or access token expired                                 |
| 404  | Resource not found                                                    |
| 409  | Attempt to create duplicate resource                                  |
| 500  | Unexpected internal error                                             |

## Endpoints

Endpoints are organized in three access levels:

- **Public**: no authentication required
- **Common**: authentication required
- **Admin**: restricted to administrators

### Public

#### Health Check

| Method | Endpoint      | Description      |
| ------ | ------------- | ---------------- |
| GET    | `/api/health` | Check API status |

#### Authentication

| Method | Endpoint             | Description                             |
| ------ | -------------------- | --------------------------------------- |
| POST   | `/api/auth/register` | User registration (requires activation) |
| POST   | `/api/auth/login`    | Login and JWT token return via cookies  |
| POST   | `/api/auth/logout`   | User logout                             |
| POST   | `/api/auth/refresh`  | Access token refresh                    |

#### Countries

| Method | Endpoint               | Description            |
| ------ | ---------------------- | ---------------------- |
| GET    | `/api/countries`       | List all countries     |
| GET    | `/api/countries/:code` | Get a specific country |

### Common

#### User

| Method | Endpoint        | Description                 |
| ------ | --------------- | --------------------------- |
| GET    | `/api/users/me` | Get authenticated user data |
| PUT    | `/api/users/me` | Update user data            |

#### Origin

| Method | Endpoint       | Description         |
| ------ | -------------- | ------------------- |
| GET    | `/api/origins` | List active origins |

#### Sequencer

| Method | Endpoint          | Description            |
| ------ | ----------------- | ---------------------- |
| GET    | `/api/sequencers` | List active sequencers |

#### Sample Source

| Method | Endpoint              | Description                |
| ------ | --------------------- | -------------------------- |
| GET    | `/api/sample-sources` | List active sample sources |

#### Laboratory

| Method | Endpoint            | Description              |
| ------ | ------------------- | ------------------------ |
| GET    | `/api/laboratories` | List active laboratories |

### Admin

Admin endpoints follow the complete CRUD pattern for **Users**, **Origins**, **Sequencers**, **Sample Sources**, and **Laboratories**:

#### User

| Method | Endpoint                                | Description              |
| ------ | --------------------------------------- | ------------------------ |
| GET    | `/api/admin/users`                      | List all users           |
| GET    | `/api/admin/users/:username`            | Get a specific user      |
| POST   | `/api/admin/users`                      | Create an activated user |
| PUT    | `/api/admin/users/:username`            | Update a user            |
| PUT    | `/api/admin/users/activation/:username` | Activate/deactivate user |
| DELETE | `/api/admin/users/:username`            | Delete a user            |

#### Origin

| Method | Endpoint                       | Description            |
| ------ | ------------------------------ | ---------------------- |
| GET    | `/api/admin/origins`           | List all origins       |
| GET    | `/api/admin/origins/:originId` | Get a specific origin  |
| PUT    | `/api/admin/origins/search`    | Search origins by name |
| POST   | `/api/admin/origins`           | Create a new origin    |
| PUT    | `/api/admin/origins/:originId` | Update an origin       |
| DELETE | `/api/admin/origins/:originId` | Delete an origin       |

#### Sequencer

| Method | Endpoint                       | Description                         |
| ------ | ------------------------------ | ----------------------------------- |
| GET    | `/api/admin/sequencers`        | List all sequencers                 |
| GET    | `/api/admin/sequencers/:id`    | Get a specific sequencer            |
| PUT    | `/api/admin/sequencers/search` | Search sequencers by brand or model |
| POST   | `/api/admin/sequencers`        | Create a new sequencer              |
| PUT    | `/api/admin/sequencers/:id`    | Update a sequencer                  |
| DELETE | `/api/admin/sequencers/:id`    | Delete a sequencer                  |

#### Sample Source

| Method | Endpoint                           | Description                            |
| ------ | ---------------------------------- | -------------------------------------- |
| GET    | `/api/admin/sample-sources`        | List all sample sources                |
| GET    | `/api/admin/sample-sources/:id`    | Get a specific sample source           |
| PUT    | `/api/admin/sample-sources/search` | Search sample sources by name or group |
| POST   | `/api/admin/sample-sources`        | Create a new sample source             |
| PUT    | `/api/admin/sample-sources/:id`    | Update a sample source                 |
| DELETE | `/api/admin/sample-sources/:id`    | Delete a sample source                 |

#### Laboratory

| Method | Endpoint                         | Description                                 |
| ------ | -------------------------------- | ------------------------------------------- |
| GET    | `/api/admin/laboratories`        | List all laboratories                       |
| GET    | `/api/admin/laboratories/:id`    | Get a specific laboratory                   |
| PUT    | `/api/admin/laboratories/search` | Search laboratories by name or abbreviation |
| POST   | `/api/admin/laboratories`        | Create a new laboratory                     |
| PUT    | `/api/admin/laboratories/:id`    | Update a laboratory                         |
| DELETE | `/api/admin/laboratories/:id`    | Delete a laboratory                         |
