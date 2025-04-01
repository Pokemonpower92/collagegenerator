package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageRepository struct {
	client *pgxpool.Pool
	ctx    context.Context
	q      *sqlc.Queries
}

type CollageReader interface {
	Get(id uuid.UUID) (*sqlc.Collage, error)
	GetAll() ([]*sqlc.Collage, error)
}

type CollageWriter interface {
	Create(sqlc.CreateCollageParams) (*sqlc.Collage, error)
}

func NewCollageRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*CollageRepository, error) {
	connString := GetConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &CollageRepository{
		client: client,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (cr *CollageRepository) Close() {
	cr.client.Close()
}

func (cr *CollageRepository) Get(id uuid.UUID) (*sqlc.Collage, error) {
	collage, err := cr.q.GetCollage(cr.ctx, id)
	if err != nil {
		return nil, err
	}
	return collage, nil
}

func (cr *CollageRepository) GetAll() ([]*sqlc.Collage, error) {
	collages, err := cr.q.ListCollages(cr.ctx)
	if err != nil {
		return nil, err
	}
	return collages, nil
}

func (cr *CollageRepository) Create(req sqlc.CreateCollageParams) (*sqlc.Collage, error) {
	collage, err := cr.q.CreateCollage(cr.ctx, req)
	if err != nil {
		return nil, err
	}
	return collage, nil
}
