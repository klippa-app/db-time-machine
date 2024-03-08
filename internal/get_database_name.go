package internal

import (
	"context"
	"fmt"
	"slices"

	"github.com/klippa-app/db-time-machine/db"
	"github.com/klippa-app/db-time-machine/internal/config"
	"github.com/klippa-app/db-time-machine/internal/hashes"
)

type MigrateFunc func(ctx context.Context, target string) error

func GenName(ctx context.Context, hash string) string {
	config := config.FromContext(ctx)
	return fmt.Sprintf("%s_%s", config.Prefix, hash[:8])
}

func GetDatabaseName(ctx context.Context, driver db.Driver, migrate MigrateFunc) (string, error) {
	hashes := hashes.FromContext(ctx)

	names := make([]string, len(hashes))
	for i := range hashes {
		names[i] = GenName(ctx, hashes[i])
	}

	current := names[0]

	databases, err := driver.List(ctx)
	if err != nil {
		return current, err
	}

	var parent string

	for _, name := range names {
		if slices.Contains(databases, name) {
			parent = name
			break
		}
	}

	if current == parent {
		return current, nil
	}

	if parent == "" {
		err = driver.Create(ctx, current)
	} else {
		err = driver.Clone(ctx, parent, current)
	}

	if err != nil {
		return current, err
	}

	if err := migrate(ctx, current); err != nil {
		if err := driver.Remove(ctx, current); err != nil {
			return current, err
		}
		return current, err
	}

	return current, nil
}
