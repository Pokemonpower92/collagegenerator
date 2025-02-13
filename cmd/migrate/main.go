package main

import (
	"log"

	"github.com/pokemonpower92/collagegenerator/internal/database"
)

func main() {
	log.Printf("Migrating database...")
	if err := database.RunMigration(); err != nil {
		panic(err)
	} else {
		log.Printf("Migration succeeded.")
	}
}
