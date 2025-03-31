package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type UserRepository struct {
	client *pgxpool.Pool
	ctx    context.Context
	q      *sqlc.Queries
}

type UserReader interface {
	Get(id uuid.UUID) (*sqlc.User, error)
	GetByResourceId(id uuid.UUID) ([]*sqlc.User, error)
	GetAll() ([]*sqlc.User, error)
}

type UserWriter interface {
	Create(sqlc.CreateUserParams) (*sqlc.User, error)
}

func NewUserRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*UserRepository, error) {
	connString := GetConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &UserRepository{
		client: client,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (ur *UserRepository) Close() {
	ur.client.Close()
}

func (ur *UserRepository) Get(id uuid.UUID) (*sqlc.User, error) {
	User, err := ur.q.GetUser(ur.ctx, id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (ur *UserRepository) GetByResourceId(userName string) ([]*sqlc.User, error) {
	User, err := ur.q.GetByUserName(ur.ctx, userName)
	if err != nil {
		return nil, err
	}
	return []*sqlc.User{User}, nil
}

func (ur *UserRepository) GetAll() ([]*sqlc.User, error) {
	User, err := ur.q.ListUsers(ur.ctx)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (ur *UserRepository) Create(
	req sqlc.CreateUserParams,
) (*sqlc.User, error) {
	imageset, err := ur.q.CreateUser(ur.ctx, req)
	if err != nil {
		return nil, err
	}
	return imageset, nil
}
