// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CollageStatus string

const (
	CollageStatusPending    CollageStatus = "pending"
	CollageStatusProcessing CollageStatus = "processing"
	CollageStatusReady      CollageStatus = "ready"
	CollageStatusFailed     CollageStatus = "failed"
)

func (e *CollageStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CollageStatus(s)
	case string:
		*e = CollageStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for CollageStatus: %T", src)
	}
	return nil
}

type NullCollageStatus struct {
	CollageStatus CollageStatus `json:"collage_status"`
	Valid         bool          `json:"valid"` // Valid is true if CollageStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCollageStatus) Scan(value interface{}) error {
	if value == nil {
		ns.CollageStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CollageStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCollageStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CollageStatus), nil
}

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
	DbID          int32         `json:"db_id"`
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	ImageSetID    uuid.UUID     `json:"image_set_id"`
	TargetImageID uuid.UUID     `json:"target_image_id"`
	CreatedAt     pgtype.Date   `json:"created_at"`
	UpdatedAt     pgtype.Date   `json:"updated_at"`
	Status        CollageStatus `json:"status"`
}

type CollageImage struct {
	DbID      int32       `json:"db_id"`
	ID        uuid.UUID   `json:"id"`
	FileName  string      `json:"file_name"`
	CollageID uuid.UUID   `json:"collage_id"`
	CreatedAt pgtype.Date `json:"created_at"`
	UpdatedAt pgtype.Date `json:"updated_at"`
}

type CollageSection struct {
	DbID      int32       `json:"db_id"`
	ID        uuid.UUID   `json:"id"`
	ImageID   uuid.UUID   `json:"image_id"`
	CollageID uuid.UUID   `json:"collage_id"`
	Section   int32       `json:"section"`
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
