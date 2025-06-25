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
	"eda-in-golang/modules/cosec"
	"eda-in-golang/modules/customers"
	"eda-in-golang/modules/depot"
	"eda-in-golang/modules/notifications"
	"eda-in-golang/modules/ordering"
	"eda-in-golang/modules/payments"
	"eda-in-golang/modules/search"
	"eda-in-golang/modules/stores"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
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
	m := monolith{svc: sys}

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(m.svc.DB())

	// init modules
	m.modules = []system.Module{
		&baskets.Module{},
		&customers.Module{},
		&depot.Module{},
		&notifications.Module{},
		&ordering.Module{},
		&payments.Module{},
		&stores.Module{},
		&search.Module{},
		&cosec.Module{},
	}

	if err = m.startupModules(); err != nil {
		return err
	}

	m.svc.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = baskets.Root(m.svc.Runner().Context(), m.svc); err != nil {
		return err
	}

	fmt.Println("started mallbots monolith application")
	defer fmt.Println("stopped mallbots monolith application")

	m.svc.Runner().Add(
		m.svc.RunWeb,
		m.svc.RunRPC,
		m.svc.RunStream,
	)

	return m.svc.Runner().Run()
}
