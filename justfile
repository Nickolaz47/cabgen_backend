set dotenv-load

BACKUP_DIR := "./backups"
DB_CONTAINER := "cabgen_backend-database-1"

default:
    @just --list

# --- Docker ---

up:
    docker compose up -d --build

down:
    docker compose down

build:
    docker compose build

restart:
    docker compose down
    docker compose up -d --build

status:
    docker compose ps

logs:
    docker compose logs -f

logs-api:
    docker compose logs -f api

logs-worker:
    docker compose logs -f worker-analysis worker-email

backup:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p "{{BACKUP_DIR}}"
    docker exec {{DB_CONTAINER}} pg_dumpall --clean -U $DB_USER -f /tmp/backup.sql
    docker cp {{DB_CONTAINER}}:/tmp/backup.sql {{BACKUP_DIR}}/backup.sql
    docker exec {{DB_CONTAINER}} rm /tmp/backup.sql
    echo "Backup saved: {{BACKUP_DIR}}/backup.sql"

restore file:
    #!/usr/bin/env bash
    set -euo pipefail
    cat {{file}} | docker exec -i {{DB_CONTAINER}} psql -U $DB_USER -d postgres
    echo "Restored from {{file}}"

# --- Podman ---

up-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml up -d --build

down-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml down

build-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml build

restart-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml down
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml up -d --build

status-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml ps

logs-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml logs -f

logs-api-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml logs -f api

logs-worker-podman:
    podman compose -f docker-compose.yaml -f docker-compose.podman.yaml logs -f worker-analysis worker-email

backup-podman:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p "{{BACKUP_DIR}}"
    podman exec {{DB_CONTAINER}} pg_dumpall --clean -U $DB_USER -f /tmp/backup.sql
    podman cp {{DB_CONTAINER}}:/tmp/backup.sql {{BACKUP_DIR}}/backup.sql
    podman exec {{DB_CONTAINER}} rm /tmp/backup.sql
    echo "Backup saved: {{BACKUP_DIR}}/backup.sql"

restore-podman file:
    #!/usr/bin/env bash
    set -euo pipefail
    cat {{file}} | podman exec -i {{DB_CONTAINER}} psql -U $DB_USER -d postgres
    echo "Restored from {{file}}"
