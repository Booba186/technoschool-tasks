package subscriber

import (
	"context"
	"encoding/json"
	"github.com/Booba186/level0/internal/cache"
	"github.com/Booba186/level0/internal/config"
	"github.com/Booba186/level0/internal/model"
	"github.com/Booba186/level0/internal/repository"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

type Subscriber struct {
	repo   *repository.OrderRepository
	cache  *cache.Cache
	reader *kafka.Reader
}

func NewSubscriber(repo *repository.OrderRepository, cache *cache.Cache, cfg *config.Config) *Subscriber {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(cfg.KafkaBrokers, ","),
		Topic:   "orders",
		GroupID: "order-group-1",
	})
	return &Subscriber{
		repo:   repo,
		cache:  cache,
		reader: r,
	}
}

func (s *Subscriber) Start(ctx context.Context) {
	log.Println("Запуск подписчика Kafka...")
	for {
		msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Ошибка чтения сообщения из Kafka: %v", err)
			continue
		}
		var order model.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Ошибка десериализации JSON: %v. Сообщение проигнорировано.", err)
			continue
		}
		if err := s.repo.SaveOrder(ctx, order); err != nil {
			log.Printf("Ошибка сохранения заказа в БД: %v", err)
			continue
		}
		s.cache.Set(order)
		log.Printf("Заказ %s успешно обработан и сохранен.", order.OrderUID)
	}
}
