package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Seed(config *pgx.ConnConfig) error {
	ctx := context.Background()
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	return nil
}
