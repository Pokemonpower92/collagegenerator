package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type TargetImageRepository struct {
	client *pgxpool.Pool
	ctx    context.Context
	q      *sqlc.Queries
}

type TargetImageReader interface {
	Get(id uuid.UUID) (*sqlc.TargetImage, error)
	GetAll() ([]*sqlc.TargetImage, error)
}

type TargetImageWriter interface {
	Create(sqlc.CreateTargetImageParams) (*sqlc.TargetImage, error)
}

func NewTargetImageRepository(
	pgConfig *config.DBConfig,
	ctx context.Context,
) (*TargetImageRepository, error) {
	connString := GetConnectionString(pgConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &TargetImageRepository{
		client: client,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (tir *TargetImageRepository) Close() {
	tir.client.Close()
}

func (tir *TargetImageRepository) Get(id uuid.UUID) (*sqlc.TargetImage, error) {
	targetImage, err := tir.q.GetTargetImage(tir.ctx, id)
	if err != nil {
		return nil, err
	}
	return targetImage, nil
}

func (tir *TargetImageRepository) GetAll() ([]*sqlc.TargetImage, error) {
	targetImages, err := tir.q.ListTargetImages(tir.ctx)
	if err != nil {
		return nil, err
	}
	return targetImages, nil
}

func (tir *TargetImageRepository) Create(
	req sqlc.CreateTargetImageParams,
) (*sqlc.TargetImage, error) {
	targetImage, err := tir.q.CreateTargetImage(tir.ctx, req)
	if err != nil {
		return nil, err
	}
	return targetImage, nil
}
