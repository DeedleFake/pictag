-- name: GetImage :one
SELECT * FROM images WHERE id = ? LIMIT 1;

-- name: CreateImage :one
INSERT INTO images (id) VALUES (?) RETURNING *;
