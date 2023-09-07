package subscriber

import (
	"log"

	"github.com/nats-io/stan.go"
)

func NatsSubscriber() {
	sc, err := stan.Connect("test-cluster", "WBsub")
	if err != nil {
		log.Fatalf("!2 %v", err)
	}
	defer sc.Close()

	subscribe, err := sc.Subscribe("WBorder", func(msg *stan.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
	}, stan.DeliverAllAvailable())
	if err != nil {
		log.Fatalf("Error subscribing to NATS Streaming channel: %v", err)
	}
	defer subscribe.Unsubscribe()

	select {}
}
