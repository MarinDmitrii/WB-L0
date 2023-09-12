package adapters

import (
	"encoding/json"

	domain "github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type OrderModel struct {
	OrderUID          string          `db:"order_uid"`
	TrackNumber       string          `db:"track_number"`
	Entry             string          `db:"entry"`
	Delivery          json.RawMessage `db:"delivery"`
	Payment           json.RawMessage `db:"payment"`
	Items             json.RawMessage `db:"items"`
	Locale            string          `db:"locale"`
	InternalSignature string          `db:"internal_signature"`
	CustomerID        string          `db:"customer_id"`
	DeliveryService   string          `db:"delivery_service"`
	Shardkey          string          `db:"shardkey"`
	SmID              int             `db:"sm_id"`
	DateCreated       string          `db:"date_created"`
	OofShard          string          `db:"oof_shard"`
}

func NewOrderModel(order domain.Order) (OrderModel, error) {
	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return OrderModel{}, err
	}

	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return OrderModel{}, err
	}

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return OrderModel{}, err
	}

	return OrderModel{
			OrderUID:          order.OrderUID,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Locale:            order.Locale,
			Delivery:          deliveryJSON,
			Payment:           paymentJSON,
			Items:             itemsJSON,
			InternalSignature: order.InternalSignature,
			CustomerID:        order.CustomerID,
			DeliveryService:   order.DeliveryService,
			Shardkey:          order.Shardkey,
			SmID:              order.SmID,
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
		},
		nil
}

func mapToDomain(order *OrderModel) domain.Order {
	var delivery domain.Delivery
	if err := json.Unmarshal([]byte(order.Delivery), &delivery); err != nil {
		return domain.Order{}
	}

	var payment domain.Payment
	if err := json.Unmarshal([]byte(order.Payment), &payment); err != nil {
		return domain.Order{}
	}

	var items []domain.Item
	if err := json.Unmarshal([]byte(order.Items), &items); err != nil {
		return domain.Order{}
	}

	return domain.Order{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
}
