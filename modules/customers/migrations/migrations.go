package migrations

import (
	"embed"
)

//go:embed *.sql
var CustomersFS embed.FS
