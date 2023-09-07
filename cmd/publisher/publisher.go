package publisher

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/MarinDmitrii/WB-L0/internal/order/domain"
	"github.com/nats-io/stan.go"
)

func NatsPublisher() {
	sc, err := stan.Connect("test-cluster", "WBpub")
	if err != nil {
		log.Fatalf("!1 %v", err)
	}
	defer sc.Close()

	for i := 0; ; i++ {
		randomOrder := RandomOrder()
		message, err := json.Marshal(randomOrder)
		if err != nil {
			log.Fatalf("Can't marshal order to json: %v\n", err)
			return
		}

		err = sc.Publish("WBorder", message)
		if err != nil {
			log.Fatalf("Can't publish message into NATS: %v\n", err)
		}

		time.Sleep(5 * time.Second)
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomOrder() *domain.Order {
	orderUID := randomString(15) + "test"
	entry := "WB" + randomString(2)
	trackNumber := entry + randomString(10)

	order := domain.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       entry,
		Delivery: domain.Delivery{
			Name:    randomString(randomInt(3, 8)) + " " + randomString(randomInt(5, 12)),
			Phone:   "+" + strconv.Itoa(randomInt(1, 999)) + strconv.Itoa(randomInt(1000000000, 9999999999)),
			Zip:     strconv.Itoa(randomInt(1000000, 9999999)),
			City:    randomString(randomInt(3, 19)),
			Address: randomString(randomInt(5, 9)) + " " + randomString(randomInt(4, 7)) + " " + strconv.Itoa(randomInt(1, 300)),
			Region:  randomString(randomInt(3, 15)),
			Email:   randomString(randomInt(4, 15)) + "@gmail.com",
		},
		Payment: domain.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       randomInt(100, 9999),
			PaymentDt:    int(time.Now().Unix()),
			Bank:         randomString(randomInt(4, 13)),
			DeliveryCost: randomInt(100, 9999),
			GoodsTotal:   randomInt(100, 999),
			CustomFee:    0,
		},
		Items: []domain.Item{
			{
				ChrtID:      randomInt(1000000, 9999999),
				TrackNumber: trackNumber,
				Price:       randomInt(100, 9999),
				Rid:         randomString(17) + "test",
				Name:        randomString(randomInt(4, 13)),
				Sale:        randomInt(0, 70),
				Size:        randomString(1),
				TotalPrice:  randomInt(100, 999),
				NmID:        randomInt(1000000, 9999999),
				Brand:       randomString(randomInt(4, 18)),
				Status:      randomInt(100, 400),
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	return &order
}
