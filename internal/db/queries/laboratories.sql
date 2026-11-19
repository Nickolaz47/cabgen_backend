-- name: GetLaboratories :many
SELECT * FROM laboratories
ORDER BY name ASC;

-- name: GetActiveLaboratories :many
SELECT * FROM laboratories
WHERE is_active = true
ORDER BY name ASC;

-- name: GetLaboratoriesByNameOrAbbreviation :many
SELECT * FROM laboratories
WHERE name ILIKE '%' || sqlc.arg(search)::text || '%'
   OR abbreviation ILIKE '%' || sqlc.arg(search)::text || '%'
ORDER BY name ASC;

-- name: GetLaboratoryByID :one
SELECT * FROM laboratories
WHERE id = $1 LIMIT 1;

-- name: GetLaboratoryDuplicate :one
SELECT * FROM laboratories
WHERE LOWER(name) = LOWER(sqlc.arg(name))
AND id != sqlc.arg(exclude_id)
LIMIT 1;

-- name: CreateLaboratory :one
INSERT INTO laboratories (
    name,
    abbreviation,
    is_active
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateLaboratory :one
UPDATE laboratories
SET
    name = $2,
    abbreviation = $3,
    is_active = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteLaboratory :exec
DELETE FROM laboratories
WHERE id = $1;
