-- name: listMigrations :many
SELECT * FROM migrations ORDER BY name ASC;

-- name: addMigration :exec
INSERT INTO migrations (name) VALUES (?);

-- name: ListImages :many
SELECT * FROM images ORDER BY created_at LIMIT ? OFFSET ?;

-- name: GetImage :one
SELECT * FROM images WHERE id = ? LIMIT 1;

-- name: CreateImage :one
INSERT INTO images (id, image_created_at) VALUES (?, ?) RETURNING *;

-- name: TagImage :one
INSERT INTO tags (name, image_id) VALUES (?, ?) RETURNING *;

-- name: SearchTags :many
SELECT name, COUNT(*) FROM tags WHERE name LIKE ? GROUP BY name ORDER BY COUNT(*) DESC, name ASC LIMIT ?;

-- name: ImagesByTags :many
SELECT images.* FROM images
JOIN tags ON tags.image_id = images.id
WHERE tags.name IN sqlc.slice('tags')
GROUP BY images.id
HAVING COUNT(DISTINCT tags.name) = CAST(@length AS INTEGER)
LIMIT ? OFFSET ?;
