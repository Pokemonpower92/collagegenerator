// Repositories are thin wrappers over the generated sqlc queries.

package repository

import (
	"fmt"

	"github.com/pokemonpower92/collagegenerator/config"
)

// Get the database connection string based on the provided
// config
func GetConnectionString(config *config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
}
