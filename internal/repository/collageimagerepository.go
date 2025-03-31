package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageImageRepository struct {
	client *pgxpool.Pool
	ctx    context.Context
	q      *sqlc.Queries
}

type CollageImageReader interface {
	Get(id uuid.UUID) (*sqlc.CollageImage, error)
	GetByResourceId(id uuid.UUID) ([]*sqlc.CollageImage, error)
	GetAll() ([]*sqlc.CollageImage, error)
}

type CollageImageWriter interface {
	Create(req uuid.UUID) (*sqlc.CollageImage, error)
}

func NewCollageImgageRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*CollageImageRepository, error) {
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
