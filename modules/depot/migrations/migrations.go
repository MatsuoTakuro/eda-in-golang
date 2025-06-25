package migrations

import (
	"embed"
)

//go:embed *.sql
var DepotFS embed.FS
