package config

import (
	"context"
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"
)

func FromContext(ctx context.Context) Config {
	value := ctx.Value(configContextKey{})
	cfg, ok := value.(Config)
	if !ok {
		panic("no config in the context")
	}
	return cfg
}

func Load(ctx context.Context, path string) (context.Context, error) {
	var config *Config
	var err error

	if path != "" {
		config, err = readConfig(path)
		if err != nil {
			return nil, err
		}
	} else {
		files := []string{".dbtm.toml", "dbtm.toml"}
		for _, file := range files {
			config, err = readConfig(file)
			if err != nil && !os.IsNotExist(err) {
				return nil, err
			}
		}
	}

	if config == nil {
		return nil, errors.New("could not find a config file")
	}

	return context.WithValue(ctx, configContextKey{}, *config), nil
}

func MergeFlags(ctx context.Context, flags *pflag.FlagSet) (context.Context, error) {
	config := FromContext(ctx)

	flags.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "uri":
			config.Connection.URI = f.Value.String()
		case "database":
			config.Connection.Database = f.Value.String()
		case "migration-directory":
			config.Migration.Directory = f.Value.String()
		case "migration-format":
			config.Migration.Format = f.Value.String()
		case "migration-command":
			config.Migration.Command = f.Value.(pflag.SliceValue).GetSlice()
		}
	})

	return context.WithValue(ctx, configContextKey{}, config), nil
}

func readConfig(path string) (*Config, error) {
	config := new(Config)

	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

type configContextKey struct{}

type Config struct {
	Prefix string

	Connection struct {
		URI      string
		Database string
	}

	Migration struct {
		Directory string
		Format    string
		Command   []string
	}
}
