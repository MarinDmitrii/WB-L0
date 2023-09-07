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
		"locale" VARCHAR,
		"internal_signature" VARCHAR,
		"customer_id" VARCHAR,
		"delivery_service" VARCHAR,
		"shardkey" VARCHAR,
		"sm_id" INTEGER,
		"date_created" TIMESTAMP,
		"oof_shard" VARCHAR
	  );
	  
	  CREATE TABLE IF NOT EXISTS "delivery" (
		"order_uid" VARCHAR PRIMARY KEY NOT NULL REFERENCES "orders"("order_uid") ON DELETE CASCADE,
		"name" VARCHAR,
		"phone" VARCHAR,
		"zip" VARCHAR,
		"city" VARCHAR,
		"address" VARCHAR,
		"region" VARCHAR,
		"email" VARCHAR
	  );
	  
	  CREATE TABLE IF NOT EXISTS "payment" (
		"order_uid" VARCHAR PRIMARY KEY NOT NULL REFERENCES "orders"("order_uid") ON DELETE CASCADE,
		"transaction" VARCHAR,
		"request_id" VARCHAR,
		"currency" VARCHAR(20),
		"provider" VARCHAR,
		"amount" INTEGER,
		"payment_dt" INTEGER,
		"bank" VARCHAR,
		"delivery_cost" INTEGER,
		"goods_total" INTEGER,
		"custom_fee" INTEGER
	  );
	  
	  CREATE TABLE IF NOT EXISTS "item" (
		"item_id" SERIAL PRIMARY KEY NOT NULL,
		"order_uid" VARCHAR NOT NULL REFERENCES "orders"("order_uid") ON DELETE CASCADE,
		"chrt_id" INTEGER,
		"track_number" VARCHAR,
		"price" INTEGER,
		"rid" VARCHAR,
		"name" VARCHAR,
		"sale" INTEGER,
		"size" VARCHAR,
		"total_price" INTEGER,
		"nm_id" INTEGER,
		"brand" VARCHAR,
		"status" INTEGER
	  );	
	`)

	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) SaveOrder(ctx context.Context, domainOrder domain.Order) (string, error) {
	transaction, err := r.db.Beginx()
	if err != nil {
		return "", err
	}

	order, delivery, payment, items, err := NewOrderModel(domainOrder)
	if err != nil {
		return "", err
	}

	queryOrder := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature,
			:customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)
		RETURNING order_uid
	`
	_, err = transaction.NamedExecContext(ctx, queryOrder, order)
	if err != nil {
		transaction.Rollback()
		return "", err
	}

	queryDelivery := `
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES (:order_uid, :name, :phone, :zip, :city, :address, :region, :email)
	`
	_, err = transaction.NamedExecContext(ctx, queryDelivery, delivery)
	if err != nil {
		transaction.Rollback()
		return "", err
	}

	queryPayment := `
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES (:order_uid, :transaction, :request_id, :currency, :provider, :amount,
			:payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)
	`
	_, err = transaction.NamedExecContext(ctx, queryPayment, payment)
	if err != nil {
		transaction.Rollback()
		return "", err
	}

	for _, item := range items {
		err := r.db.QueryRowContext(ctx, "SELECT nextval('item_item_id_seq')").Scan(&item.ItemID)
		if err != nil {
			return "", err
		}

		queryItem := `
			INSERT INTO item (item_id, order_uid, chrt_id, track_number, price,
				rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES (:item_id, :order_uid, :chrt_id, :track_number, :price,
				:rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status)
		`
		_, err = transaction.NamedExecContext(ctx, queryItem, item)

		if err != nil {
			transaction.Rollback()
			return "", err
		}
	}

	if err := transaction.Commit(); err != nil {
		return "", err
	}

	return order.OrderUID, nil
}

// func (r *PostgresOrderRepository) GetOrderByID(ctx context.Context, orderUID string) (domain.Order, error) {
// 	var order *OrderModel
// 	queryOrder := `SELECT * FROM orders WHERE order_uid = $1`
// 	err := r.db.GetContext(ctx, &order, queryOrder, orderUID)
// 	if err != nil {
// 		fmt.Println("1")
// 		return domain.Order{}, err
// 	}

// 	var delivery *DeliveryModel
// 	queryDelivery := `SELECT * FROM delivery WHERE order_uid = $1`
// 	err = r.db.GetContext(ctx, &delivery, queryDelivery, orderUID)
// 	if err != nil {
// 		fmt.Println("2")
// 		return domain.Order{}, err
// 	}

// 	var payment *PaymentModel
// 	queryPayment := `SELECT * FROM payment WHERE order_uid = '$1'`
// 	err = r.db.GetContext(ctx, &payment, queryPayment, orderUID)
// 	if err != nil {
// 		fmt.Println("3")
// 		return domain.Order{}, err
// 	}

// 	var items []*ItemModel
// 	queryItems := `SELECT * FROM item WHERE order_uid = '$1'`
// 	fmt.Println(orderUID)
// 	err = r.db.SelectContext(ctx, &items, queryItems, orderUID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("4")
// 			// Обработка случая, когда нет данных
// 		} else {
// 			fmt.Println("5")
// 			return domain.Order{}, err
// 		}
// 	}

// 	return mapToDomain(order, delivery, payment, items), nil
// }
