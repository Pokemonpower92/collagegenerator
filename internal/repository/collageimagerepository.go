package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageImageRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
	ctx    context.Context
	q      *sqlc.Queries
}

func NewCollageImgageRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*CollageImageRepository, error) {
	logger := log.New(
		log.Writer(),
		"CollageImageRepository: ",
		log.LstdFlags,
	)
	connString := GetConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &CollageImageRepository{
		client: client,
		logger: logger,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (cir *CollageImageRepository) Close() {
	cir.client.Close()
}

func (cir *CollageImageRepository) Get(id uuid.UUID) (*sqlc.CollageImage, error) {
	collageImage, err := cir.q.GetCollageImage(cir.ctx, id)
	if err != nil {
		return nil, err
	}
	return collageImage, nil
}

func (cir *CollageImageRepository) GetByResourceId(id uuid.UUID) ([]*sqlc.CollageImage, error) {
	collageImage, err := cir.q.GetByCollageId(cir.ctx, id)
	if err != nil {
		return nil, err
	}
	return []*sqlc.CollageImage{collageImage}, nil
}

func (cir *CollageImageRepository) GetAll() ([]*sqlc.CollageImage, error) {
	collageImage, err := cir.q.ListCollageImages(cir.ctx)
	if err != nil {
		return nil, err
	}
	return collageImage, nil
}

func (cir *CollageImageRepository) Create(
	req uuid.UUID,
) (*sqlc.CollageImage, error) {
	imageset, err := cir.q.CreateCollageImage(cir.ctx, req)
	if err != nil {
		return nil, err
	}
	return imageset, nil
}

func (cir *CollageImageRepository) Update(
	id uuid.UUID,
	req uuid.UUID,
) (*sqlc.CollageImage, error) {
	cir.logger.Printf("Update not implemented for CollageImage")
	return nil, errors.New("Not implemented")
}

func (cir *CollageImageRepository) Delete(id uuid.UUID) error {
	cir.logger.Printf("Delete not implemented for CollageImage")
	return errors.New("Not implemented")
}
