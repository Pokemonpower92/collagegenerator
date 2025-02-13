package main

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/database"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
)

func main() {
	log.Printf("Migrating database...")
	config.LoadEnvironmentVariables()
	postgresConfig := config.NewPostgresConfig()
	connString := repository.GetConnectionString(postgresConfig)
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	if err := database.RunMigration(config); err != nil {
		panic(err)
	} else {
		log.Printf("Migration succeeded.")
	}
}
