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
	"eda-in-golang/modules/customers"
	"eda-in-golang/modules/customers/migrations"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("customers exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(s.DB())
	if err = s.MigrateDB(migrations.CustomersFS); err != nil {
		return err
	}
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = customers.Root(s.Runner().Context(), s); err != nil {
		return err
	}

	fmt.Println("started customers service")
	defer fmt.Println("stopped customers service")

	s.Runner().Add(
		s.RunWeb,
		s.RunRPC,
		s.RunStream,
	)

	// go func() {
	// 	for {
	// 		var mem runtime.MemStats
	// 		runtime.ReadMemStats(&mem)
	// 		m.logger.Debug().Msgf("Alloc = %v  TotalAlloc = %v  Sys = %v  NumGC = %v", mem.Alloc/1024, mem.TotalAlloc/1024, mem.Sys/1024, mem.NumGC)
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	return s.Runner().Run()
}
