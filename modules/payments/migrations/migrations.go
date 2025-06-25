package migrations

import (
	"embed"
)

//go:embed *.sql
var PaymentsFS embed.FS
