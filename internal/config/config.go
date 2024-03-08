package config

import (
	"context"
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

func FromContext(ctx context.Context) Config {
	value := ctx.Value(configContextKey{})
	cfg, ok := value.(Config)
	if !ok {
		panic("no config in the context")
	}
	return cfg
}

func Attach(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configContextKey{}, cfg)
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

	return Attach(ctx, *config), nil
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
	Prefix string `json:",omitempty"`

	Connection ConnectionConfig
	Migration  MigrationConfig
}

type ConnectionConfig struct {
	URI      string `json:",omitempty"`
	Database string `json:",omitempty"`
}

type MigrationConfig struct {
	Directory string `json:",omitempty"`
	Format    string `json:",omitempty"`
	Command   string `json:",omitempty"`
}
