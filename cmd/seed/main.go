package main

import (
	"log"

	"github.com/pokemonpower92/collagegenerator/internal/database"
)

func main() {
	log.Printf("Seeding database...")
	database.Seed()
	log.Printf("Seeding succeeded.")
}
