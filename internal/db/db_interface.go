package utils

import (
	"context"
	"database/sql"
)

type db interface {
	getConnection(ctx context.Context) *sql.DB
	List(ctx context.Context) ([]string, error)
	Clone(ctx context.Context, clonedDBName string, newDBName string) error
	Remove(ctx context.Context, DBName string) error
	Prune(ctx context.Context) error
}
