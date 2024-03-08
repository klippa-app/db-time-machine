package hashes

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/klippa-app/db-time-machine/internal/config"
)

type Hashes []string

/*
	Calculate determines the hashes migrations and attaches them to the context.

*
  - Each file name is hashed along with the hash of the previous file name;
  - This means that if any file is added, removed, or renamed it will break
  - the hash chain and invalidate any database that was created after.
  - The file contents is then hashed with the chained file name hash;
    This means that if the contents of a file changes, it will invalidate
    the database only for that migration, but not any subsiquent databases.

*
*
*/
func Calculate(ctx context.Context) (context.Context, error) {
	cfg := config.FromContext(ctx)
	nameRegex, err := regexp.Compile(cfg.Migration.Format)
	if err != nil {
		return nil, err
	}

	migrations, err := filepath.Abs(cfg.Migration.Directory)
	if err != nil {
		return nil, err
	}

	chainSum := [16]byte{}
	hashes := []string{}

	filepath.Walk(migrations, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !nameRegex.MatchString(info.Name()) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		chainSum = md5.Sum(append([]byte(info.Name()), chainSum[:]...))
		dataSum := md5.Sum(append(data, chainSum[:]...))
		hash := hex.EncodeToString(dataSum[:])

		hashes = append([]string{hash}, hashes...)

		return nil
	})

	return Attach(ctx, hashes), nil
}

type hashesKey struct{}

func Attach(ctx context.Context, hashes Hashes) context.Context {
	return context.WithValue(ctx, hashesKey{}, hashes)
}

func FromContext(ctx context.Context) Hashes {
	hashes, ok := ctx.Value(hashesKey{}).(Hashes)
	if !ok || hashes == nil {
		panic(errors.New("no hashes on the context"))
	}

	return hashes
}
