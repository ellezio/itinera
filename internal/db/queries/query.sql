-- name: GetResource :one
SELECT sqlc.embed(r), sqlc.embed(s) FROM resources r
JOIN statuses s ON r.status_id = s.id
WHERE r.id = ? LIMIT 1;

-- name: GetResources :many
SELECT sqlc.embed(r), sqlc.embed(s) FROM resources r
JOIN statuses s ON r.status_id = s.id;

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

-- name: GetStatuses :many
SELECT * FROM statuses;

-- name: GetResourceTags :many
SELECT t.*
FROM resource_tags rt
JOIN tags t ON rt.tag_id = t.id
WHERE rt.resource_id = ?;

-- name: GetResourcesTags :many
SELECT rt.resource_id, sqlc.embed(t)
FROM resource_tags rt
JOIN tags t ON rt.tag_id = t.id
WHERE rt.resource_id IN (sqlc.slice('resources'));

-- name: GetTags :many
SELECT * FROM tags;

-- name: SetTag :exec
INSERT INTO resource_tags (
  resource_id, tag_id
) VALUES (
  ?, ?
);

-- name: SetStatus :exec
UPDATE resources
SET status_id = ?
WHERE id = ?;

-- name: ClearTags :exec
DELETE FROM resource_tags
WHERE resource_id = ?;

-- name: GetResourcesNotes :many
SELECT * FROM notes
WHERE entity_id IN (sqlc.slice('resources')) AND entity_type=?;

-- name: GetNotes :many
SELECT * FROM notes
WHERE entity_id=? AND entity_type=?;

-- name: GetNote :one
SELECT * FROM notes
WHERE id=? LIMIT 1;

-- name: CreateNote :one
INSERT INTO notes (
  title, content, entity_id, entity_type
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id=?;

-- name: UpdateNote :one
UPDATE notes SET 
title=?,
content=?
WHERE id=?
RETURNING *;

-- name: GetStatus :one
SELECT * FROM statuses WHERE id=?;

-- name: CreateTag :one
INSERT INTO tags (
 name, color
) VALUES (
  ?, ?
)
RETURNING *;

-- name: CreateStatus :one
INSERT INTO statuses (
 name, color
) VALUES (
  ?, ?
)
RETURNING *;


-- name: GetTag :one
SELECT * FROM tags
WHERE id=?;

-- name: UpdateStatus :one
UPDATE statuses SET
name=?,
color=?
WHERE id=?
RETURNING *;

-- name: UpdateTag :one
UPDATE tags SET
name=?,
color=?
WHERE id=?
RETURNING *;

-- name: DeleteStatus :exec
DELETE FROM statuses WHERE id=?;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id=?;
