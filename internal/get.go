package internal

import (
	"context"
	"fmt"

	"github.com/klippa-app/db-time-machine/db"
	"github.com/klippa-app/db-time-machine/internal/config"
	"github.com/klippa-app/db-time-machine/internal/hashes"
)

type MigrateFunc func(ctx context.Context, target string) error

func GenName(ctx context.Context, hash string) string {
	config := config.FromContext(ctx)
	return fmt.Sprintf("%s_%s", config.Prefix, hash)
}

func GetHashList(ctx context.Context) []string {
	// hash list should be lastest migration first
	config := config.FromContext(ctx)
	_ = config
	return nil
}

func NearestParent(ctx context.Context, hashes []string) string {
	config := config.FromContext(ctx)
	_ = config
	return ""
}

func Get(ctx context.Context, driver db.Driver, migrate MigrateFunc) (string, error) {
	hashes := hashes.FromContext(ctx)
	currentName := GenName(ctx, hashes[0])
	parentName := NearestParent(ctx, hashes)

	if currentName == parentName {
		return currentName, nil
	}

	err := driver.Clone(ctx, parentName, currentName)
	if err != nil {
		return currentName, err
	}

	err = migrate(ctx, currentName)
	if err != nil {
		// Nuke the failed db?
		return currentName, err
	}

	return currentName, nil
}
