-- name: GetResource :one
SELECT * FROM resources
WHERE id = ? LIMIT 1;

-- name: GetResources :many
SELECT * FROM resources;

-- name: CreateResource :one
INSERT INTO resources (
  title, source, source_type, status_id
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateResource :one
UPDATE resources SET
title = ?,
source = ?,
source_type = ?,
status_id = ?
WHERE id = ?
RETURNING *;

-- name: DeleteResource :exec
DELETE FROM resources
WHERE id = ?;
