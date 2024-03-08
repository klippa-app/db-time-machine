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

	last := [16]byte{}
	hashes := []string{}

	filepath.Walk(migrations, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !nameRegex.MatchString(info.Name()) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sum := md5.Sum(append(data, last[:]...))
		hash := hex.EncodeToString(sum[:])

		hashes = append([]string{hash}, hashes...)
		last = sum

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
