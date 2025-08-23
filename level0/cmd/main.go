package main

import (
	"context"
	"github.com/Booba186/level0/internal/cache"
	"github.com/Booba186/level0/internal/config"
	"github.com/Booba186/level0/internal/handler"
	"github.com/Booba186/level0/internal/repository"
	"github.com/Booba186/level0/internal/subscriber"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()

	dbPool, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer dbPool.Close()
	log.Println("Успешное подключение к базе данных!")

	orderRepo := repository.NewOrderRepository(dbPool)
	orderCache := cache.New()
	log.Println("Кэш инициализирован.")

	orders, err := orderRepo.GetAllOrders(context.Background())
	if err != nil {
		log.Fatalf("Ошибка при загрузке заказов из БД: %v", err)
	}
	for _, order := range orders {
		orderCache.Set(order)
	}
	log.Printf("Кэш восстановлен. Загружено %d заказов.\n", len(orders))

	sub := subscriber.NewSubscriber(orderRepo, orderCache, cfg)
	go sub.Start(context.Background())

	//http
	h := handler.NewHandler(orderCache)
	router := chi.NewRouter()
	router.Get("/order/{uid}", h.GetOrderByUID)
	router.Handle("/*", http.FileServer(http.Dir("./web")))

	log.Println("Сервис запущен на порту :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %v", err)
	}
}
