//go:build tools
// +build tools

package main

import (
	_ "github.com/klippa-app/db-time-machine"
	_ "github.com/rubenv/sql-migrate/sql-migrate"
)
