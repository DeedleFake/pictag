-- name: listMigrations :many
SELECT * FROM migrations ORDER BY name ASC;

-- name: addMigration :exec
INSERT INTO migrations (name) VALUES (?);

-- name: GetImage :one
SELECT * FROM images WHERE id = ? LIMIT 1;

-- name: CreateImage :one
INSERT INTO images (id, image_created_at) VALUES (?, ?) RETURNING *;

-- name: TagImage :one
INSERT INTO tags (name, image_id) VALUES (?, ?) RETURNING *;

-- name: SearchTags :many
SELECT name, COUNT(*) FROM tags WHERE name LIKE ? GROUP BY name ORDER BY COUNT(*) DESC, name ASC LIMIT ?;
