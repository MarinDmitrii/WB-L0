package adapters

import (
	"time"

	domain "github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type OrderModel struct {
	OrderUID          string    `db:"order_uid"`
	TrackNumber       string    `db:"track_number"`
	Entry             string    `db:"entry"`
	Locale            string    `db:"locale"`
	InternalSignature string    `db:"internal_signature"`
	CustomerID        string    `db:"customer_id"`
	DeliveryService   string    `db:"delivery_service"`
	Shardkey          string    `db:"shardkey"`
	SmID              int       `db:"sm_id"`
	DateCreated       time.Time `db:"date_created"`
	OofShard          string    `db:"oof_shard"`
}

type DeliveryModel struct {
	OrderUID string `db:"order_uid"`
	Name     string `db:"name"`
	Phone    string `db:"phone"`
	Zip      string `db:"zip"`
	City     string `db:"city"`
	Address  string `db:"address"`
	Region   string `db:"region"`
	Email    string `db:"email"`
}

type PaymentModel struct {
	OrderUID     string `db:"order_uid"`
	Transaction  string `db:"transaction"`
	RequestID    string `db:"request_id"`
	Currency     string `db:"currency"`
	Provider     string `db:"provider"`
	Amount       int    `db:"amount"`
	PaymentDt    int    `db:"payment_dt"`
	Bank         string `db:"bank"`
	DeliveryCost int    `db:"delivery_cost"`
	GoodsTotal   int    `db:"goods_total"`
	CustomFee    int    `db:"custom_fee"`
}

type ItemModel struct {
	ItemID      int    `db:"item_id"`
	OrderUID    string `db:"order_uid"`
	ChrtID      int    `db:"chrt_id"`
	TrackNumber string `db:"track_number"`
	Price       int    `db:"price"`
	Rid         string `db:"rid"`
	Name        string `db:"name"`
	Sale        int    `db:"sale"`
	Size        string `db:"size"`
	TotalPrice  int    `db:"total_price"`
	NmID        int    `db:"nm_id"`
	Brand       string `db:"brand"`
	Status      int    `db:"status"`
}

func NewOrderModel(order domain.Order) (OrderModel, DeliveryModel, PaymentModel, []ItemModel, error) {
	var itemModels []ItemModel
	for _, item := range order.Items {
		itemModel := ItemModel{
			OrderUID:    order.OrderUID,
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Rid:         item.Rid,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		}
		itemModels = append(itemModels, itemModel)
	}

	return OrderModel{
			OrderUID:          order.OrderUID,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Locale:            order.Locale,
			InternalSignature: order.InternalSignature,
			CustomerID:        order.CustomerID,
			DeliveryService:   order.DeliveryService,
			Shardkey:          order.Shardkey,
			SmID:              order.SmID,
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
		},
		DeliveryModel{
			OrderUID: order.OrderUID,
			Name:     order.Delivery.Name,
			Phone:    order.Delivery.Phone,
			Zip:      order.Delivery.Zip,
			City:     order.Delivery.City,
			Address:  order.Delivery.Address,
			Region:   order.Delivery.Region,
			Email:    order.Delivery.Email,
		},
		PaymentModel{
			OrderUID:     order.OrderUID,
			Transaction:  order.Payment.Transaction,
			RequestID:    order.Payment.RequestID,
			Currency:     order.Payment.Currency,
			Provider:     order.Payment.Provider,
			Amount:       order.Payment.Amount,
			PaymentDt:    order.Payment.PaymentDt,
			Bank:         order.Payment.Bank,
			DeliveryCost: order.Payment.DeliveryCost,
			GoodsTotal:   order.Payment.GoodsTotal,
			CustomFee:    order.Payment.CustomFee,
		},
		itemModels,
		nil
}

// func mapToDomain(order *OrderModel, delivery *DeliveryModel, payment *PaymentModel, items []*ItemModel) domain.Order {
// 	itemss := make([]domain.Item, len(items))
// 	for i, itemModel := range items {
// 		item := domain.Item{
// 			ChrtID:      itemModel.ChrtID,
// 			TrackNumber: itemModel.TrackNumber,
// 			Price:       itemModel.Price,
// 			Rid:         itemModel.Rid,
// 			Name:        itemModel.Name,
// 			Sale:        itemModel.Sale,
// 			Size:        itemModel.Size,
// 			TotalPrice:  itemModel.TotalPrice,
// 			NmID:        itemModel.NmID,
// 			Brand:       itemModel.Brand,
// 			Status:      itemModel.Status,
// 		}
// 		itemss[i] = item
// 	}

// 	return domain.Order{
// 		OrderUID:    order.OrderUID,
// 		TrackNumber: order.TrackNumber,
// 		Entry:       order.Entry,
// 		Delivery: domain.Delivery{
// 			Name:    delivery.Name,
// 			Phone:   delivery.Phone,
// 			Zip:     delivery.Zip,
// 			City:    delivery.City,
// 			Address: delivery.Address,
// 			Region:  delivery.Region,
// 			Email:   delivery.Email,
// 		},
// 		Payment: domain.Payment{
// 			Transaction:  payment.Transaction,
// 			RequestID:    payment.RequestID,
// 			Currency:     payment.Currency,
// 			Provider:     payment.Provider,
// 			Amount:       payment.Amount,
// 			PaymentDt:    payment.PaymentDt,
// 			Bank:         payment.Bank,
// 			DeliveryCost: payment.DeliveryCost,
// 			GoodsTotal:   payment.GoodsTotal,
// 			CustomFee:    payment.CustomFee,
// 		},
// 		Items:             itemss,
// 		Locale:            order.Locale,
// 		InternalSignature: order.InternalSignature,
// 		CustomerID:        order.CustomerID,
// 		DeliveryService:   order.DeliveryService,
// 		Shardkey:          order.Shardkey,
// 		SmID:              order.SmID,
// 		DateCreated:       order.DateCreated,
// 		OofShard:          order.OofShard,
// 	}
// }
