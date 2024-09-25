package pkg

import (
	"app/internal/config"
	"app/pkg/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

func NewPsqlClient(ctx context.Context, cfg *config.Config) (pool *pgxpool.Pool, err error) {
	queryConnection := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?%s",
		cfg.PostgresCfg.User,
		cfg.PostgresCfg.Password,
		cfg.PostgresCfg.Host,
		cfg.PostgresCfg.Port,
		cfg.PostgresCfg.DbName,
	)

	//TODO: DoWithTries использвать для подключения эту библиотеку https://github.com/kamilsk/retry
	err = utils.DoWithTries(func() error {

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, queryConnection)
		if err != nil {
			return err
		}
		return nil

	}, 4, 5*time.Second)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
