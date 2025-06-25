package migrations

import (
	"embed"
)

// MonoFS provides embedded SQL migration files for the mono-repo structure.
//
//go:embed *.sql
var MonoFS embed.FS
