package ports

import (
	"net/http"

	"github.com/MarinDmitrii/WB-L0/internal/order/builder"
	"github.com/MarinDmitrii/WB-L0/internal/order/domain"
	"github.com/MarinDmitrii/WB-L0/internal/order/usecase"
	"github.com/labstack/echo/v4"
)

type HttpOrderHandler struct {
	app *builder.Application
}

func NewHttpOrderHandler(app *builder.Application) HttpOrderHandler {
	return HttpOrderHandler{app: app}
}

func (h HttpOrderHandler) SaveOrder(ctx echo.Context) error {
	request := &PostOrder{}
	if err := ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	saveOrder := &usecase.SaveOrder{
		Order: domain.Order{
			OrderUID:    request.OrderUID,
			TrackNumber: request.TrackNumber,
			Entry:       request.Entry,
			Delivery: domain.Delivery{
				Name:    request.Delivery.Name,
				Phone:   request.Delivery.Phone,
				Zip:     request.Delivery.Zip,
				City:    request.Delivery.City,
				Address: request.Delivery.Address,
				Region:  request.Delivery.Region,
				Email:   request.Delivery.Email,
			},
			Payment: domain.Payment{
				Transaction:  request.Payment.Transaction,
				RequestID:    request.Payment.RequestID,
				Currency:     request.Payment.Currency,
				Provider:     request.Payment.Provider,
				Amount:       request.Payment.Amount,
				PaymentDt:    request.Payment.PaymentDt,
				Bank:         request.Payment.Bank,
				DeliveryCost: request.Payment.DeliveryCost,
				GoodsTotal:   request.Payment.GoodsTotal,
				CustomFee:    request.Payment.CustomFee,
			},
			Items:             make([]domain.Item, len(request.Items)),
			Locale:            request.Locale,
			InternalSignature: request.InternalSignature,
			CustomerID:        request.CustomerID,
			DeliveryService:   request.DeliveryService,
			Shardkey:          request.Shardkey,
			SmID:              request.SmID,
			DateCreated:       request.DateCreated,
			OofShard:          request.OofShard,
		},
	}

	for i, p := range request.Items {
		saveOrder.Order.Items[i] = domain.Item{
			ChrtID:      p.ChrtID,
			TrackNumber: p.TrackNumber,
			Price:       p.Price,
			Rid:         p.Rid,
			Name:        p.Name,
			Sale:        p.Sale,
			Size:        p.Size,
			TotalPrice:  p.TotalPrice,
			NmID:        p.NmID,
			Brand:       p.Brand,
			Status:      p.Status,
		}
	}

	orderUID, err := h.app.SaveOrder.Execute(ctx.Request().Context(), saveOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, struct {
		OrderUID string `json:"order_uid"`
	}{OrderUID: orderUID})
}

func (h HttpOrderHandler) GetOrderByID(ctx echo.Context, orderUID string) error {
	orders, err := h.app.GetOrderByID.Execute(ctx.Request().Context(), orderUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	response := h.mapToResponse(orders)

	return ctx.JSON(http.StatusOK, response)
}

func (h HttpOrderHandler) mapToResponse(order domain.Order) Order {
	var items []Item
	for _, itemm := range order.Items {
		item := Item{
			ChrtID:      itemm.ChrtID,
			TrackNumber: itemm.TrackNumber,
			Price:       itemm.Price,
			Rid:         itemm.Rid,
			Name:        itemm.Name,
			Sale:        itemm.Sale,
			Size:        itemm.Size,
			TotalPrice:  itemm.TotalPrice,
			NmID:        itemm.NmID,
			Brand:       itemm.Brand,
			Status:      itemm.Status,
		}
		items = append(items, item)
	}

	return Order{
		CustomerID:  order.CustomerID,
		DateCreated: order.DateCreated,
		Delivery: Delivery{
			Address: order.Delivery.Address,
			City:    order.Delivery.City,
			Email:   order.Delivery.Email,
			Name:    order.Delivery.Name,
			Phone:   order.Delivery.Phone,
			Region:  order.Delivery.Region,
			Zip:     order.Delivery.Zip,
		},
		DeliveryService:   order.DeliveryService,
		Entry:             order.Entry,
		InternalSignature: order.InternalSignature,
		Items:             items,
		Locale:            order.Locale,
		OofShard:          order.OofShard,
		OrderUID:          order.OrderUID,
		Payment: Payment{
			Amount:       order.Payment.Amount,
			Bank:         order.Payment.Bank,
			Currency:     order.Payment.Currency,
			CustomFee:    order.Payment.CustomFee,
			DeliveryCost: order.Payment.DeliveryCost,
			GoodsTotal:   order.Payment.GoodsTotal,
			PaymentDt:    order.Payment.PaymentDt,
			Provider:     order.Payment.Provider,
			RequestID:    order.Payment.RequestID,
			Transaction:  order.Payment.Transaction,
		},
		Shardkey:    order.Shardkey,
		SmID:        order.SmID,
		TrackNumber: order.TrackNumber,
	}

}

func CustomRegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST("/orders", wrapper.CreateOrder)
	router.GET("/orders/:order_uid", wrapper.GetOrderByID)
}
