-- name: GetTargetImage :one
SELECT * FROM target_images
WHERE id = $1 LIMIT 1;

-- name: ListTargetImages :many
SELECT * FROM target_images
ORDER BY name;

-- name: CreateTargetImage :one
INSERT INTO target_images (
  id, name, description, created_at, updated_at
) VALUES (
  $1, $2, $3, NOW(), NOW() 
)
RETURNING *;

