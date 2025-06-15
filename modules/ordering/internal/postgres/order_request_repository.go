package postgres

import (
	"context"
	pg "eda-in-golang/internal/postgres"
	"eda-in-golang/modules/ordering/internal/domain"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type orderRequestRepository struct {
	tableName string
	db        pg.DB
}

var _ domain.OrderRequestRepository = (*orderRequestRepository)(nil)

func NewOrderRequestRepository(tableName string, db pg.DB) *orderRequestRepository {
	return &orderRequestRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r orderRequestRepository) FindOrInsert(ctx context.Context,
	idemKey string, request domain.OrderRequest, command any,
) (orderID string, inserted bool, err error) {

	const query = `
		WITH attempt AS (
			INSERT INTO %s (
				id, idempotency_key, type, command_payload
			) VALUES ($1, $2, $3, $4)
			ON CONFLICT (idempotency_key) DO NOTHING
			RETURNING id, TRUE AS inserted
		)
		SELECT id, inserted FROM attempt
		UNION
		SELECT id, FALSE AS inserted
		FROM %s
		WHERE idempotency_key = $2
		LIMIT 1;
	`

	payload, err := json.Marshal(command)
	if err != nil {
		return "", false, fmt.Errorf("marshal command payload: %w", err)
	}

	newID := uuid.New().String()
	row := r.db.QueryRowContext(ctx, fmt.Sprintf(query, r.tableName, r.tableName),
		newID, idemKey, string(request), payload,
	)

	err = row.Scan(&orderID, &inserted)
	if err != nil {
		return "", false, fmt.Errorf("find-or-insert order request: %w", err)
	}

	return orderID, inserted, nil
}

func (r orderRequestRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
