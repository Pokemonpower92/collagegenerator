package metadata

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/client"
	"github.com/pokemonpower92/collagegenerator/internal/service"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

func Start() {
	rmq, err := client.NewRabbitMQClient(config.NewRMQConfig(slog.Default()))
	if err != nil {
		panic(err)
	}
	rmq.StartReceiving(
		context.Background(),
		config.METADATA_QUEUE(),
		func(message string) error {
			var collage sqlc.Collage
			err := json.Unmarshal([]byte(message), &collage)
			if err != nil {
				slog.Default().Error("Error unmarshaling collage", "error", err)
				return err
			}
			service.CreateCollageMetaData(&collage, slog.Default())
			return nil
		},
	)
}
