//go:build integration || database

// NOTE:
// This file is excluded from code navigation by default
// unless the matching build tags are passed.
//
// If you're using VS Code, add the following to `.vscode/settings.json`:
// "gopls": {
//     "buildFlags": [
//         "-tags=integration database",
//     ],
// }

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"eda-in-golang/internal/logger/log"
	"eda-in-golang/migrations"
	"eda-in-golang/modules/baskets/internal/domain"
	"slices"
)

type productCacheSuite struct {
	container testcontainers.Container
	db        *sql.DB
	mock      *domain.MockProductClient
	repo      ProductCacheRepository
	suite.Suite
}

// go test -tags=integration,database eda-in-golang/modules/baskets/internal/postgres -v
func TestProductCacheRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode: skipping")
	}
	suite.Run(t, &productCacheSuite{})
}

// TestProductCacheRepository(t *testing.T) : entry point for the test suite
//    ├─ SetupSuite()        : run once before any tests
//    │   ├─ SetupTest()     : run before each test
//    │   │   └─ TestXxx()   : actual test functions (run per test case)
//    │   └─ TearDownTest()  : run after each test
//    └─ TearDownSuite()     : run once after all tests

func (s *productCacheSuite) SetupSuite() {
	var err error

	ctx := context.Background()
	initDir, err := filepath.Abs("./../../../../docker/database")
	if err != nil {
		s.T().Fatal(err)
	}
	entries, err := os.ReadDir(initDir)
	if err != nil {
		s.T().Fatalf("failed to read init dir: %v", err)
	}
	initDBfilenames := []string{
		"001_create_monolith_db.sh",
		"002_create_pact_db.sh",
		// Ignore other table init files that are similar to another migration files as migrations.FS.
	}
	initDBfiles := make([]testcontainers.ContainerFile, 0, len(initDBfilenames))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sh") {
			continue
		}
		if !slices.Contains(initDBfilenames, entry.Name()) {
			continue
		}

		initDBfiles = append(initDBfiles, testcontainers.ContainerFile{
			HostFilePath:      filepath.Join(initDir, entry.Name()),
			ContainerFilePath: filepath.Join("/docker-entrypoint-initdb.d", entry.Name()),
			FileMode:          0755, // executable
		})
	}
	const dbURL = "postgres://mallbots_user:mallbots_pass@%s:%s/mallbots?sslmode=disable"
	s.container, err = testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "postgres:12-alpine",
				Hostname:     "postgres",
				ExposedPorts: []string{"5432/tcp"},
				Env: map[string]string{
					"POSTGRES_PASSWORD": "itsasecret",
				},
				Files: initDBfiles,
				WaitingFor: wait.ForSQL("5432/tcp", "pgx",
					func(host string, port nat.Port) string {
						return fmt.Sprintf(dbURL, host, port.Port())
					},
				).WithStartupTimeout(5 * time.Second),
			},
			Started: true,
		})
	if err != nil {
		s.T().Fatal(err)
	}

	endpoint, err := s.container.Endpoint(ctx, "")
	if err != nil {
		s.T().Fatal(err)
	}

	s.db, err = sql.Open("pgx", fmt.Sprintf("postgres://mallbots_user:mallbots_pass@%s/mallbots?sslmode=disable", endpoint))
	if err != nil {
		s.T().Fatal(err)
	}

	goose.SetLogger(&log.SilentLogger{})
	goose.SetBaseFS(migrations.MonoFS)
	if err := goose.SetDialect("postgres"); err != nil {
		s.T().Fatal(err)
	}
	if err := goose.Up(s.db, "."); err != nil {
		s.T().Fatal(err)
	}
}

func (s *productCacheSuite) TearDownSuite() {
	err := s.db.Close()
	if err != nil {
		s.T().Fatal(err)
	}
	if err := s.container.Terminate(context.Background()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *productCacheSuite) SetupTest() {
	// Instead of doing this in TearDownTest, clean up leftover data before each test
	// to ensure a consistent starting state. This avoids test order dependency and improves reliability.
	// TearDownTest may be skipped if the test crashes or is interrupted in a debugger.
	_, err := s.db.ExecContext(context.Background(), "TRUNCATE baskets.products_cache")
	if err != nil {
		s.T().Fatal(err)
	}

	s.mock = domain.NewMockProductClient(s.T())
	s.repo = NewProductCacheRepository("baskets.products_cache", s.db, s.mock)
}

func (s *productCacheSuite) TestProductCacheRepository_Add() {
	// Act
	s.NoError(s.repo.Add(context.Background(), "product-id", "store-id", "product-name", 10.00))
	row := s.db.QueryRow("SELECT name FROM baskets.products_cache WHERE id = $1", "product-id")

	// Assert
	if s.NoError(row.Err()) {
		var name string
		s.NoError(row.Scan(&name))
		s.Equal("product-name", name)
	}
}

func (s *productCacheSuite) TestProductCacheRepository_AddDupe() {
	// Act
	s.NoError(s.repo.Add(context.Background(), "product-id", "store-id", "product-name", 10.00))
	s.NoError(s.repo.Add(context.Background(), "product-id", "store-id", "dupe-product-name", 10.00))
	row := s.db.QueryRow("SELECT name FROM baskets.products_cache WHERE id = $1", "product-id")

	// Assert
	if s.NoError(row.Err()) {
		var name string
		s.NoError(row.Scan(&name))
		s.Equal("product-name", name)
	}
}

func (s *productCacheSuite) TestProductCacheRepository_Rebrand() {
	// Arrange
	_, err := s.db.Exec("INSERT INTO baskets.products_cache (id, store_id, name, price) VALUES ('product-id', 'store-id', 'product-name', 10.00)")
	s.NoError(err)

	// Act
	s.NoError(s.repo.Rebrand(context.Background(), "product-id", "new-product-name"))

	// Assert
	row := s.db.QueryRow("SELECT name FROM baskets.products_cache WHERE id = $1", "product-id")
	if s.NoError(row.Err()) {
		var name string
		s.NoError(row.Scan(&name))
		s.Equal("new-product-name", name)
	}
}

func (s *productCacheSuite) TestProductCacheRepository_UpdatePrice() {
	// Arrange
	_, err := s.db.Exec("INSERT INTO baskets.products_cache (id, store_id, name, price) VALUES ('product-id', 'store-id', 'product-name', 10.00)")
	s.NoError(err)

	// Act
	s.NoError(s.repo.UpdatePrice(context.Background(), "product-id", 2.00))
	row := s.db.QueryRow("SELECT price FROM baskets.products_cache WHERE id = $1", "product-id")

	// Assert
	if s.NoError(row.Err()) {
		var price float64
		s.NoError(row.Scan(&price))
		s.Equal(12.00, price)
	}
}
func (s *productCacheSuite) TestProductCacheRepository_Remove() {
	// Arrange
	_, err := s.db.Exec("INSERT INTO baskets.products_cache (id, store_id, name, price) VALUES ('product-id', 'store-id', 'product-name', 10.00)")
	s.NoError(err)

	// Act
	s.NoError(s.repo.Remove(context.Background(), "product-id"))
	row := s.db.QueryRow("SELECT name FROM baskets.products_cache WHERE id = $1", "product-id")

	// Assert
	if s.NoError(row.Err()) {
		var name string
		s.Error(row.Scan(&name))
	}
}

func (s *productCacheSuite) TestProductCacheRepository_Find() {
	// Arrange
	_, err := s.db.Exec("INSERT INTO baskets.products_cache (id, store_id, name, price) VALUES ('product-id', 'store-id', 'product-name', 10.00)")
	s.NoError(err)

	// Act
	product, err := s.repo.Find(context.Background(), "product-id")

	// Assert
	if s.NoError(err) {
		s.NotNil(product)
		s.Equal("product-name", product.Name)
	}
}

func (s *productCacheSuite) TestProductCacheRepository_FindFromFallback() {
	// Arrange
	s.mock.On("Find", mock.Anything, "product-id").Return(&domain.Product{
		ID:      "product-id",
		StoreID: "store-id",
		Name:    "product-name",
		Price:   10.00,
	}, nil)

	// Act
	product, err := s.repo.Find(context.Background(), "product-id")

	// Assert
	if s.NoError(err) {
		s.NotNil(product)
		s.Equal("product-name", product.Name)
	}
}
