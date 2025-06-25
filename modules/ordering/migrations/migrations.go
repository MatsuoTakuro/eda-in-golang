package migrations

import (
	"embed"
)

//go:embed *.sql
var OrderingFS embed.FS
