package dialects

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Postgres struct {
	ConnectionString string
	connection       *sql.DB
}

func (p Postgres) getConnection(ctx context.Context) *sql.DB {
	if p.connection == nil {
		db, err := sql.Open("postgres", p.ConnectionString)
		if err != nil {
			log.Fatal(err)
		}

		p.connection = db
	}

	return p.connection
}

func (p Postgres) List(ctx context.Context) ([]string, error) {
	prefix := "test%"
	// we need to add the % to the prefix ourselves and not in the query, it might be possible to do so there but this is the easiest solution right now.
	// the connection string we should make ourselves from the config we can get from the context or elsewhere

	rows, err := p.getConnection(ctx).Query("SELECT datname FROM pg_database WHERE datname LIKE $1", prefix)
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

func (p Postgres) Clone(ctx context.Context, clonedDBName string, newDBName string) error {
	_, err := p.getConnection(ctx).Exec(
		fmt.Sprintf(
			"CREATE DATABASE %s WITH TEMPLATE %s OWNER %s",
			pq.QuoteIdentifier(newDBName),
			pq.QuoteIdentifier(clonedDBName),
			"dochorizon"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) Remove(ctx context.Context, DBName string) error {
	_, err := p.getConnection(ctx).Exec(
		fmt.Sprintf("DROP DATABASE %s", DBName),
	)
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) Prune(ctx context.Context) error {
	// instead of regex we might want to do it based off off either last x amount of databases we have or maybe not used within last x days.
	r, err := regexp.Compile(fmt.Sprintf("^%s", "test"))
	if err != nil {
		return err
	}

	databases, err := p.List(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < len(databases); i++ {
		if r.Match([]byte(databases[i])) {
			p.Remove(ctx, databases[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
