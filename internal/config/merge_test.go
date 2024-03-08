package config

import (
	"encoding/json"
	"slices"
	"testing"
)

func emptyConfig() Config {
	return Config{}
}

func filledConfig(value string) Config {
	return Config{
		Prefix: value,
		Connection: ConnectionConfig{
			URI:      value,
			Database: value,
		},
		Migration: MigrationConfig{
			Directory: value,
			Format:    value,
			Command:   value,
		},
	}
}

func equal(a, b any) bool {
	adat, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	bdat, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	return slices.Equal(adat, bdat)
}

func TestLazyCopy(t *testing.T) {
	var src, dst Config

	src = emptyConfig()
	dst = filledConfig("a")
	lazyCopy(src, &dst)

	if !equal(dst, filledConfig("a")) {
		t.Fatalf("empty config should not overwrite filled config")
	}

	src = filledConfig("b")
	dst = filledConfig("a")
	lazyCopy(src, &dst)

	if !equal(src, dst) {
		t.Fatalf("filled src should overwrite dest")
	}

	src = Config{
		Connection: ConnectionConfig{
			URI: "c",
		},
	}
	dst = filledConfig("a")
	lazyCopy(src, &dst)

	expected := filledConfig("a")
	expected.Connection.URI = "c"

	if !equal(dst, expected) {
		t.Fatalf("partial src should only overwrite src filled fields")
	}

	t.Logf("%+v", dst)
}
