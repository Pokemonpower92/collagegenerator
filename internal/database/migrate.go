package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	migrate "github.com/rubenv/sql-migrate"
)

func RunMigration() error {
	config.LoadEnvironmentVariables()
	ctx := context.Background()
	postgresConfig := config.NewPostgresConfig()
	connString := repository.GetConnectionString(postgresConfig)

	// Create a connection config
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		return err
	}

	// Create a connection
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// Set up the migration source
	migrations := &migrate.FileMigrationSource{
		Dir: "internal/sqlc/migrations/postgres",
	}

	// Create a *sql.DB instance for migrations
	db := stdlib.OpenDB(*config)
	defer db.Close()

	// Run the migrations
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
