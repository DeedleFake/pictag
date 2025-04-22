-- name: GetImage :one
SELECT * FROM images WHERE id = ? LIMIT 1;

-- name: CreateImage :one
INSERT INTO images (id, image_created_at) VALUES (?, ?) RETURNING *;

-- name: TagImage :one
INSERT INTO tags (name, image_id) VALUES (?, ?) RETURNING *;
