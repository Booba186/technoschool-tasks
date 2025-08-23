package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9292"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	log.Println("Отправка тестового сообщения в Kafka...")

	msg := `{
		"order_uid": "b563feb7b2b84b6test",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
			"name": "Test Testov", "phone": "+9720000000", "zip": "2639809",
			"city": "Kiryat Mozkin", "address": "Ploshad Mira 15", "region": "Kraiot", "email": "test@gmail.com"
		},
		"payment": {
			"transaction": "b563feb7b2b84b6test", "currency": "USD", "provider": "wbpay",
			"amount": 1817, "payment_dt": 1637907727, "bank": "alpha", "delivery_cost": 1500
		},
		"items": [
			{
				"chrt_id": 9934930, "track_number": "WBILMTESTTRACK", "price": 453, "rid": "ab4219087a764ae0btest",
				"name": "Mascaras", "sale": 30, "size": "0", "total_price": 317, "nm_id": 2389222, "brand": "Vivienne Sabo"
			}
		],
		"locale": "en",
		"customer_id": "test",
		"date_created": "2021-11-26T06:22:19Z"
	}`

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(msg),
		},
	)

	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	log.Println("Сообщение успешно отправлено!")
}
