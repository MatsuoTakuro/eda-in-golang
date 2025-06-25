package migrations

import (
	"embed"
)

//go:embed *.sql
var BasketsFS embed.FS
