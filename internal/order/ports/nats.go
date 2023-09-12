package ports

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MarinDmitrii/WB-L0/internal/order/builder"
	"github.com/MarinDmitrii/WB-L0/internal/order/domain"
	"github.com/nats-io/stan.go"
)

type NatsOrderHandler struct {
	app *builder.Application
}

func NewNatsOrderHandler(app *builder.Application) NatsOrderHandler {
	return NatsOrderHandler{app: app}
}

func (h *NatsOrderHandler) NatsSubscriber(ctx context.Context, doneCh <-chan struct{}) error {
	subscribe, err := h.app.StanConnect.Subscribe("WBorder", func(msg *stan.Msg) {
		order := domain.Order{}
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			return
		}
		log.Printf("Received message, OrderUID = %v\n", order.OrderUID)

		_, err := h.app.SaveOrder.Execute(ctx, order)
		if err != nil {
			log.Printf("Error saving order: %v", err)
			return
		}
	}, stan.DeliverAllAvailable())

	if err != nil {
		log.Fatalf("Error subscribing to NATS Streaming channel: %v", err)
	}

	defer subscribe.Unsubscribe()

	// select {}

	for {
		select {
		case <-ctx.Done():
			log.Println("NatsSubscriber has been terminated")
			return nil
		case <-doneCh:
			log.Println("NatsSubscriber has been told to exit")
			return nil
		}
	}
}
