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

type TargetImageRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
	ctx    context.Context
	q      *sqlc.Queries
}

func NewTargetImageRepository(
	pgConfig *config.DBConfig,
	ctx context.Context,
) (*TargetImageRepository, error) {
	logger := log.New(
		log.Writer(),
		"TargetImageRepository: ",
		log.LstdFlags,
	)
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
		logger: logger,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (tir *TargetImageRepository) Close() {
	tir.client.Close()
}

func (tir *TargetImageRepository) Get(id uuid.UUID) (*sqlc.TargetImage, error) {
	tir.logger.Printf("Get not implemented")
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

func (tir *TargetImageRepository) Update(
	id uuid.UUID,
	req sqlc.CreateTargetImageParams,
) (*sqlc.TargetImage, error) {
	tir.logger.Printf("Update not implemented")
	return nil, errors.New("Update not implemented for target images")
}

func (tir *TargetImageRepository) Delete(id uuid.UUID) error {
	tir.logger.Printf("Delete not implemented")
	return errors.New("Delete not implemented for target images")
}
