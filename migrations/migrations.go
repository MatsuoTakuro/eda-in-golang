package migrations

import (
	"embed"
)

// FS provides embedded SQL migration files for Goose to run without external file access.
//
//go:embed *.sql
var FS embed.FS
