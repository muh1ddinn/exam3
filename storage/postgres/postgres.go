package postgres

import (
	"context"
	"exam3/config"
	"exam3/pkg/logger"
	"exam3/storage"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

type Store struct {
	Pool   *pgxpool.Pool
	logger logger.ILogger
}

func New(ctx context.Context, cfg config.Config, logger logger.ILogger) (storage.IStorage, error) {
	url := fmt.Sprintf(`host=%s port=%v user=%s password=%s database=%s sslmode=disable`,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDatabase)

	pgPoolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pgPoolConfig.MaxConns = 100
	pgPoolConfig.MaxConnLifetime = time.Hour

	newPool, err := pgxpool.NewWithConfig(context.Background(), pgPoolConfig)

	if err != nil {
		fmt.Println("error while connecting to db ", err.Error())
		return nil, err
	}
	return Store{
		Pool: newPool,
	}, nil
}
func (s Store) CloseDB() {
	s.Pool.Close()
}

func (s Store) Users() storage.IUsersStorage {
	Newuser := NewUser(s.Pool, s.logger)

	return &Newuser
}
