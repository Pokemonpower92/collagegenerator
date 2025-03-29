// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AverageColor struct {
	DbID       int32       `json:"db_id"`
	ID         uuid.UUID   `json:"id"`
	ImagesetID uuid.UUID   `json:"imageset_id"`
	FileName   string      `json:"file_name"`
	R          int32       `json:"r"`
	G          int32       `json:"g"`
	B          int32       `json:"b"`
	A          int32       `json:"a"`
	CreatedAt  pgtype.Date `json:"created_at"`
	UpdatedAt  pgtype.Date `json:"updated_at"`
}

type Collage struct {
	DbID          int32       `json:"db_id"`
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	ImageSetID    uuid.UUID   `json:"image_set_id"`
	TargetImageID uuid.UUID   `json:"target_image_id"`
	CreatedAt     pgtype.Date `json:"created_at"`
	UpdatedAt     pgtype.Date `json:"updated_at"`
}

type CollageImage struct {
	DbID      int32       `json:"db_id"`
	ID        uuid.UUID   `json:"id"`
	CollageID uuid.UUID   `json:"collage_id"`
	CreatedAt pgtype.Date `json:"created_at"`
	UpdatedAt pgtype.Date `json:"updated_at"`
}

type ImageSet struct {
	DbID        int32       `json:"db_id"`
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CreatedAt   pgtype.Date `json:"created_at"`
	UpdatedAt   pgtype.Date `json:"updated_at"`
}

type TargetImage struct {
	DbID        int32       `json:"db_id"`
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CreatedAt   pgtype.Date `json:"created_at"`
	UpdatedAt   pgtype.Date `json:"updated_at"`
}

type User struct {
	DbID        int32       `json:"db_id"`
	ID          uuid.UUID   `json:"id"`
	UserName    string      `json:"user_name"`
	Permissions []byte      `json:"permissions"`
	CreatedAt   pgtype.Date `json:"created_at"`
	UpdatedAt   pgtype.Date `json:"updated_at"`
}
