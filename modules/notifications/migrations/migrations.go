package migrations

import (
	"embed"
)

//go:embed *.sql
var NotificationsFS embed.FS
