-- name: GetByCollageId :many
SELECT * FROM collage_sections
WHERE collage_id = $1;

-- name: CreateCollageSection :one
INSERT INTO collage_sections (
  id, image_id, collage_id, section, created_at, updated_at
) VALUES (
  uuid_generate_v4(), $1, $2, $3, NOW(), NOW() 
)
RETURNING *;