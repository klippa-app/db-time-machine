package dbtm

import (
	"context"

	"github.com/klippa-app/db-time-machine/db"
	"github.com/klippa-app/db-time-machine/internal"
	"github.com/klippa-app/db-time-machine/internal/config"
	"github.com/klippa-app/db-time-machine/internal/hashes"
)

func TimeTravel(driver db.Driver, migrateFn MigrateFunc, cfg Config) (string, error) {
	ctx := context.Background()
	ctx, err := config.Load(ctx, cfg.ConfigFile)
	if err != nil {
		return "", err
	}

	ctx, err = config.Merge(ctx, cfg.Config)
	if err != nil {
		return "", err
	}

	ctx, err = hashes.Calculate(ctx)
	if err != nil {
		panic(err)
	}

	ctx = db.AttachContext(ctx, driver)

	return internal.TimeTravel(ctx, internal.MigrateFunc(migrateFn))
}

type Config struct {
	config.Config
	ConfigFile string
}

type MigrateFunc internal.MigrateFunc
