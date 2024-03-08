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
	Create(ctx context.Context, target string) error
	Remove(ctx context.Context, target string) error
	PruneList(ctx context.Context) ([]string, error)
	Prune(ctx context.Context) error
}

type driverContextKey struct{}

func FromContext(ctx context.Context) Driver {
	value := ctx.Value(driverContextKey{})
	driver, ok := value.(Driver)
	if !ok {
		panic("no driver in the context")
	}
	return driver
}

func AttachContext(ctx context.Context, driver Driver) context.Context {
	return context.WithValue(ctx, driverContextKey{}, driver)
}
