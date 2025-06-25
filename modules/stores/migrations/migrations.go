package migrations

import (
	"embed"
)

//go:embed *.sql
var StoresFS embed.FS
