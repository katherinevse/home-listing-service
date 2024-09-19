package user

import (
	"context"
	"github.com/jackc/pgconn"
)

type DBPool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}
