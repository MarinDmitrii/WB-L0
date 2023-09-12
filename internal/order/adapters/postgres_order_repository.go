package adapters

import (
	"context"

	domain "github.com/MarinDmitrii/WB-L0/internal/order/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresOrderRepository struct {
	db *sqlx.DB
}

func NewPostgresOrderRepository(db *sqlx.DB) *PostgresOrderRepository {
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS "orders" (
		"order_uid" VARCHAR PRIMARY KEY NOT NULL,
		"track_number" VARCHAR,
		"entry" VARCHAR,
		"delivery" jsonb,
		"payment" jsonb,
		"items" jsonb,
		"locale" VARCHAR,
		"internal_signature" VARCHAR,
		"customer_id" VARCHAR,
		"delivery_service" VARCHAR,
		"shardkey" VARCHAR,
		"sm_id" INTEGER,
		"date_created" TIMESTAMP,
		"oof_shard" VARCHAR
	  );	
	`)

	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(ctx context.Context, domainOrder domain.Order) (string, error) {
	transaction, err := r.db.Beginx()
	if err != nil {
		return "", err
	}

	order, err := NewOrderModel(domainOrder)
	if err != nil {
		return "", err
	}

	queryOrder := `
		INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES (:order_uid, :track_number, :entry, :delivery, :payment, :items, :locale, :internal_signature,
			:customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)
		RETURNING order_uid
	`
	_, err = transaction.NamedExecContext(ctx, queryOrder, order)
	if err != nil {
		transaction.Rollback()
		return "", err
	}

	if err := transaction.Commit(); err != nil {
		return "", err
	}

	return order.OrderUID, nil
}

func (r *PostgresOrderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	var orders []OrderModel
	queryAll := `SELECT * FROM orders`

	err := r.db.SelectContext(ctx, &orders, queryAll)
	if err != nil {
		return nil, err
	}

	var domainOrders []domain.Order
	for _, order := range orders {
		domainOrders = append(domainOrders, mapToDomain(&order))
	}

	return domainOrders, nil
}

func (r *PostgresOrderRepository) GetOrderByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	var order *OrderModel
	queryOrder := `SELECT * FROM orders WHERE order_uid = $1`
	err := r.db.GetContext(ctx, &order, queryOrder, orderUID)
	if err != nil {
		return domain.Order{}, err
	}

	return mapToDomain(order), nil
}
