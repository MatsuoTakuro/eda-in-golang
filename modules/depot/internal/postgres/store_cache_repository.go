package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"eda-in-golang/internal/postgres"
	"eda-in-golang/modules/depot/internal/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"
)

type StoreCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.StoreRepository
}

var _ domain.StoreCacheRepository = (*StoreCacheRepository)(nil)

func NewStoreCacheRepository(tableName string, db postgres.DB, fallback domain.StoreRepository) StoreCacheRepository {
	return StoreCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r StoreCacheRepository) Add(ctx context.Context, storeID, name, location string) error {
	const query = "INSERT INTO %s (id, name, location) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name, location)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil
			}
		}
	}

	return err
}

func (r StoreCacheRepository) Rename(ctx context.Context, storeID, name string) error {
	const query = "UPDATE %s SET name = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)

	return err
}

func (r StoreCacheRepository) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	const query = "SELECT name, location FROM %s WHERE id = $1 LIMIT 1"

	store := &domain.Store{
		ID: storeID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.Name, &store.Location)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning store")
		}
		store, err = r.fallback.Find(ctx, storeID)
		if err != nil {
			return nil, errors.Wrap(err, "store fallback failed")
		}
		// attempt to add it to the cache
		return store, r.Add(ctx, store.ID, store.Name, store.Location)
	}

	return store, nil
}

func (r StoreCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
