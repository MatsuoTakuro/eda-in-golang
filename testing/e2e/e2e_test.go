//go:build e2e

package e2e

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/go-openapi/runtime/client"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stackus/dotenv"
)

func TestEndToEnd(t *testing.T) {
	sc := &suiteContext{
		transport: client.New("localhost:8080", "/", nil),
	}

	var err error
	sc.db, err = sql.Open("pgx", postgresDSN(t))
	if err != nil {
		t.Fatal(err)
	}

	suite := newTestSuite(suiteConfig{
		paths: []string{
			"features/baskets",
			"features/customers",
			"features/kiosk",
			"features/orders",
			"features/stores",
		},
		featureCtxs: []featureContext{
			sc,
			newCustomersContext(sc),
			newStoresContext(sc),
		},
	})

	if status := suite.Run(); status != 0 {
		t.Error("end to end feature test failed with status:", status)
	}
}

func postgresDSN(t *testing.T) string {
	err := dotenv.Load()
	if err != nil {
		t.Logf("warning: failed to load .env: %v", err)
	}

	// The user must have permissions to important operations like TRUNCATE.
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres" // super user by default
	}

	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = "itsasecret"
	}

	return fmt.Sprintf("postgres://%s:%s@localhost:5432/mallbots?sslmode=disable", user, pass)
}
