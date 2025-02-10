package service

import (
	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/config"
)

type CollageMetaData struct {
	Resolution config.ResolutionConfig `json:"resolution"`
	SectionMap []uuid.UUID             `json:"section_map"`
}
