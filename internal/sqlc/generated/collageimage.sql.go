// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: collageimage.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createCollageImage = `-- name: CreateCollageImage :one
INSERT INTO collage_images (
  id, collage_id, created_at, updated_at
) VALUES (
  $1, $2, NOW(), NOW() 
)
RETURNING db_id, id, collage_id, created_at, updated_at
`

type CreateCollageImageParams struct {
	ID        uuid.UUID `json:"id"`
	CollageID uuid.UUID `json:"collage_id"`
}

func (q *Queries) CreateCollageImage(ctx context.Context, arg CreateCollageImageParams) (*CollageImage, error) {
	row := q.db.QueryRow(ctx, createCollageImage, arg.ID, arg.CollageID)
	var i CollageImage
	err := row.Scan(
		&i.DbID,
		&i.ID,
		&i.CollageID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getByCollageId = `-- name: GetByCollageId :one
SELECT db_id, id, collage_id, created_at, updated_at FROM collage_images
WHERE collage_id = $1 LIMIT 1
`

func (q *Queries) GetByCollageId(ctx context.Context, collageID uuid.UUID) (*CollageImage, error) {
	row := q.db.QueryRow(ctx, getByCollageId, collageID)
	var i CollageImage
	err := row.Scan(
		&i.DbID,
		&i.ID,
		&i.CollageID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getCollageImage = `-- name: GetCollageImage :one
SELECT db_id, id, collage_id, created_at, updated_at FROM collage_images
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCollageImage(ctx context.Context, id uuid.UUID) (*CollageImage, error) {
	row := q.db.QueryRow(ctx, getCollageImage, id)
	var i CollageImage
	err := row.Scan(
		&i.DbID,
		&i.ID,
		&i.CollageID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listCollageImages = `-- name: ListCollageImages :many
SELECT db_id, id, collage_id, created_at, updated_at FROM collage_images
ORDER BY updated_at
`

func (q *Queries) ListCollageImages(ctx context.Context) ([]*CollageImage, error) {
	rows, err := q.db.Query(ctx, listCollageImages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*CollageImage
	for rows.Next() {
		var i CollageImage
		if err := rows.Scan(
			&i.DbID,
			&i.ID,
			&i.CollageID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
