package tests

import (
	"context"
	"database/sql"
)

type db_test struct {
}

func DBTest() db_test {
	return db_test{}
}

func (dt db_test) URI(ctx context.Context) string {
	return "test_URI"
}

func (dt db_test) Connection(ctx context.Context) (*sql.DB, error) {
	return nil, nil
}

func (dt db_test) List(ctx context.Context) ([]string, error) {
	return_array := []string{"test_925e267e", "test_4984f751", "test_6c2f4616"}
	return return_array, nil
}

func (dt db_test) Clone(ctx context.Context, source string, target string) error {
	return nil
}

func (dt db_test) Create(ctx context.Context, target string) error {
	return nil
}

func (dt db_test) Remove(ctx context.Context, target string) error {
	return nil
}

func (dt db_test) PruneList(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

func (dt db_test) Prune(ctx context.Context) error {
	return nil
}
