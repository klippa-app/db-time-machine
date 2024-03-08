package dialect

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/klippa-app/db-time-machine/internal"
	"github.com/klippa-app/db-time-machine/internal/config"
	"github.com/klippa-app/db-time-machine/internal/hashes"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type postgres struct {
	connection *sql.DB
}

func Postgres() postgres {
	return postgres{}
}

func (p postgres) URI(ctx context.Context) string {
	cfg := config.FromContext(ctx)
	database := cfg.Connection.Database
	return strings.Replace(cfg.Connection.URI, "{}", database, 1)
}

func (p postgres) Connection(ctx context.Context) (*sql.DB, error) {
	if p.connection == nil {
		db, err := sql.Open("postgres", p.URI(ctx))
		if err != nil {
			return nil, err
		}

		p.connection = db
	}
	return p.connection, nil
}

func (p postgres) List(ctx context.Context) ([]string, error) {
	cfg := config.FromContext(ctx)

	conn, err := p.Connection(ctx)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT datname FROM pg_database WHERE datname LIKE '%s%%'", cfg.Prefix)
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string

	for rows.Next() {
		var datname string
		if err := rows.Scan(&datname); err != nil {
			return nil, err
		}
		databases = append(databases, datname)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

func (p postgres) Clone(ctx context.Context, source string, target string) error {
	conn, err := p.Connection(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(
		fmt.Sprintf(
			"CREATE DATABASE %s WITH TEMPLATE %s",
			pq.QuoteIdentifier(target),
			pq.QuoteIdentifier(source)),
	)
	if err != nil {
		return err
	}

	return nil
}

func (p postgres) Create(ctx context.Context, target string) error {
	conn, err := p.Connection(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", target))
	if err != nil {
		return err
	}

	return nil
}

func (p postgres) Remove(ctx context.Context, target string) error {
	conn, err := p.Connection(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(
		fmt.Sprintf("DROP DATABASE %s", target),
	)
	if err != nil {
		return err
	}

	return nil
}

func (p postgres) PruneList(ctx context.Context) ([]string, error) {
	hashList := hashes.FromContext(ctx)

	dbList, err := p.List(ctx)
	if err != nil {
		return nil, err
	}

	var hashNames []string
	for _, hash := range hashList {
		hashNames = append(hashNames, internal.GenName(ctx, hash))
	}

	var dbToPrune []string
	for _, db := range dbList {
		if !slices.Contains(hashNames, db) {
			dbToPrune = append(dbToPrune, db)
		}
	}

	return dbToPrune, nil
}

func (p postgres) Prune(ctx context.Context) error {
	toPrune, err := p.PruneList(ctx)
	if err != nil {
		return err
	}

	for _, db := range toPrune {
		err = p.Remove(ctx, db)
		if err != nil {
			return err
		}
	}

	return nil
}
