package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

	"eda-in-golang/internal/config"
	"eda-in-golang/internal/system"
	"eda-in-golang/internal/web"
	"eda-in-golang/modules/baskets"
	"eda-in-golang/modules/baskets/migrations"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("baskets exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	sys, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(sys.DB())
	if err = sys.MigrateDB(migrations.BasketsFS); err != nil {
		return err
	}
	sys.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = baskets.Root(sys.Runner().Context(), sys); err != nil {
		return err
	}

	fmt.Println("started baskets service")
	defer fmt.Println("stopped baskets service")

	sys.Runner().Add(
		sys.RunWeb,
		sys.RunRPC,
		sys.RunStream,
	)

	return sys.Runner().Run()
}
