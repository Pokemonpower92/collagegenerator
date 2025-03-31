package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type AverageColorRepository struct {
	client *pgxpool.Pool
	ctx    context.Context
	q      *sqlc.Queries
}

type AverageColorReader interface {
	Get(id uuid.UUID) (*sqlc.AverageColor, error)
	GetByResourceId(id uuid.UUID) ([]*sqlc.AverageColor, error)
	GetAll() ([]*sqlc.AverageColor, error)
}

type AverageColorWriter interface {
	Create(sqlc.CreateAverageColorParams) (*sqlc.AverageColor, error)
}

func NewAverageColorRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*AverageColorRepository, error) {
	connString := GetConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &AverageColorRepository{
		client: client,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (acr *AverageColorRepository) Close() {
	acr.client.Close()
}

func (acr *AverageColorRepository) Get(id uuid.UUID) (*sqlc.AverageColor, error) {
	averageColor, err := acr.q.GetAverageColor(acr.ctx, id)
	if err != nil {
		return nil, err
	}
	return averageColor, nil
}

func (acr *AverageColorRepository) GetByResourceId(id uuid.UUID) ([]*sqlc.AverageColor, error) {
	averageColors, err := acr.q.GetByImageSetId(acr.ctx, id)
	if err != nil {
		return nil, err
	}
	return averageColors, nil
}

func (acr *AverageColorRepository) GetAll() ([]*sqlc.AverageColor, error) {
	averageColors, err := acr.q.ListAverageColors(acr.ctx)
	if err != nil {
		return nil, err
	}
	return averageColors, nil
}

func (acr *AverageColorRepository) Create(
	req sqlc.CreateAverageColorParams,
) (*sqlc.AverageColor, error) {
	averageColor, err := acr.q.CreateAverageColor(acr.ctx, req)
	if err != nil {
		return nil, err
	}
	return averageColor, nil
}
