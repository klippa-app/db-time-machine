package db

import (
	"context"
	"database/sql"
)

type Driver interface {
	URI(ctx context.Context) string
	Connection(ctx context.Context) (*sql.DB, error)
	List(ctx context.Context) ([]string, error)
	Clone(ctx context.Context, clonedDBName string, newDBName string) error
	Remove(ctx context.Context, DBName string) error
	Prune(ctx context.Context) error
}
