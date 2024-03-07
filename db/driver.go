package db

import (
	"context"
	"database/sql"
)

type Driver interface {
	URI(ctx context.Context) string
	Connection(ctx context.Context) (*sql.DB, error)
	List(ctx context.Context) ([]string, error)
	Clone(ctx context.Context, source string, target string) error
	Remove(ctx context.Context, target string) error
	Prune(ctx context.Context) error
}
