package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type ImageSetRepository struct {
	client *pgxpool.Pool
	ctx    context.Context
	q      *sqlc.Queries
}

type ImageSetReader interface {
	Get(id uuid.UUID) (*sqlc.ImageSet, error)
	GetAll() ([]*sqlc.ImageSet, error)
}

type ImageSetWriter interface {
	Create(sqlc.CreateImageSetParams) (*sqlc.ImageSet, error)
}

func NewImageSetRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*ImageSetRepository, error) {
	connString := GetConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &ImageSetRepository{
		client: client,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (isr *ImageSetRepository) Close() {
	isr.client.Close()
}

func (isr *ImageSetRepository) Get(id uuid.UUID) (*sqlc.ImageSet, error) {
	imageSet, err := isr.q.GetImageSet(isr.ctx, id)
	if err != nil {
		return nil, err
	}
	return imageSet, nil
}

func (isr *ImageSetRepository) GetAll() ([]*sqlc.ImageSet, error) {
	imageSets, err := isr.q.ListImageSets(isr.ctx)
	if err != nil {
		return nil, err
	}
	return imageSets, nil
}

func (isr *ImageSetRepository) Create(req sqlc.CreateImageSetParams) (*sqlc.ImageSet, error) {
	imageset, err := isr.q.CreateImageSet(isr.ctx, req)
	if err != nil {
		return nil, err
	}
	return imageset, nil
}
