package config

import (
	"context"
	"encoding/json"

	"github.com/spf13/pflag"
)

func lazyCopy(src Config, dst *Config) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, dst)
	if err != nil {
		return err
	}

	return nil
}

func Merge(ctx context.Context, src Config) (context.Context, error) {
	dst := FromContext(ctx)

	err := lazyCopy(src, &dst)
	if err != nil {
		return ctx, err
	}

	return Attach(ctx, dst), nil
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
			config.Migration.Command = f.Value.String()
		}
	})

	return context.WithValue(ctx, configContextKey{}, config), nil
}
