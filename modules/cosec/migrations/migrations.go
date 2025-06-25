package migrations

import (
	"embed"
)

//go:embed *.sql
var CosecFS embed.FS
