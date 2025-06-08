package thumbnail

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
		config.THUMBNAIL_QUEUE(),
		func(message string) error {
			var collageImage sqlc.CollageImage
			err := json.Unmarshal([]byte(message), &collageImage)
			if err != nil {
				slog.Default().Error("Error unmarshaling CollageImage", "error", err)
				return err
			}
			service.GenerateCollage(&collageImage, slog.Default())
			return nil
		},
	)
}
