package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"

	"eda-in-golang/modules/baskets/internal/domain"
)

type productCacheRepository struct {
	tableName string
	db        *sql.DB
	fallback  domain.ProductClient
}

var _ domain.ProductCacheRepository = (*productCacheRepository)(nil)

func NewProductCacheRepository(tableName string, db *sql.DB, fallback domain.ProductClient) productCacheRepository {
	return productCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r productCacheRepository) Add(ctx context.Context, productID, storeID, name string, price float64) error {
	const query = `INSERT INTO %s (id, store_id, name, price) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, storeID, name, price)
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

func (r productCacheRepository) Rebrand(ctx context.Context, productID, name string) error {
	const query = `UPDATE %s SET name = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, name)

	return err
}

func (r productCacheRepository) UpdatePrice(ctx context.Context, productID string, delta float64) error {
	const query = `UPDATE %s SET price = price + $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, delta)

	return err
}

func (r productCacheRepository) Remove(ctx context.Context, productID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID)

	return err
}

func (r productCacheRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	const query = `SELECT store_id, name, price FROM %s WHERE id = $1 LIMIT 1`

	product := &domain.Product{
		ID: productID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(&product.StoreID, &product.Name, &product.Price)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning product")
		}
		log.Printf("Product %s not found in cache, falling back to service...", productID)
		product, err = r.fallback.Find(ctx, productID)
		if err != nil {
			return nil, errors.Wrap(err, "product fallback failed")
		}
		// attempt to add it to the cache
		return product, r.Add(ctx, product.ID, product.StoreID, product.Name, product.Price)
	}

	return product, nil
}

func (r productCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
