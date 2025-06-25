package migrations

import (
	"embed"
)

//go:embed *.sql
var SearchFS embed.FS
