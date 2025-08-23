package repository

import (
	"context"
	"fmt"
	"github.com/Booba186/level0/internal/config"
	"github.com/Booba186/level0/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepository — структура для работы с бд
type OrderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository — конструктор для репозитория
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// NewPostgresDB подключение к бд
func NewPostgresDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDBName)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("не удалось пропинговать базу данных: %w", err)
	}

	return pool, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order model.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	ordersQuery := `INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, date_created) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := tx.Exec(ctx, ordersQuery, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.CustomerID, order.DateCreated); err != nil {
		return fmt.Errorf("не удалось вставить в orders: %w", err)
	}

	deliveryQuery := `INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := tx.Exec(ctx, deliveryQuery, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email); err != nil {
		return fmt.Errorf("не удалось вставить в delivery: %w", err)
	}

	paymentQuery := `INSERT INTO payment (transaction_uid, order_uid, currency, provider, amount, payment_dt, bank, delivery_cost) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := tx.Exec(ctx, paymentQuery, order.Payment.Transaction, order.OrderUID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost); err != nil {
		return fmt.Errorf("не удалось вставить в payment: %w", err)
	}

	itemsQuery := `INSERT INTO items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	for _, item := range order.Items {
		if _, err := tx.Exec(ctx, itemsQuery, item.ChrtID, order.OrderUID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand); err != nil {
			return fmt.Errorf("не удалось вставить item с chrt_id %d: %w", item.ChrtID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	query := `SELECT order_uid, track_number, entry, locale, customer_id, date_created FROM orders`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса всех заказов: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.CustomerID, &order.DateCreated); err != nil {
			return nil, fmt.Errorf("ошибка сканирования заказа: %w", err)
		}

		delivery, err := r.getDeliveryByOrderUID(ctx, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения delivery для заказа %s: %w", order.OrderUID, err)
		}
		order.Delivery = delivery

		payment, err := r.getPaymentByOrderUID(ctx, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения payment для заказа %s: %w", order.OrderUID, err)
		}
		order.Payment = payment

		items, err := r.getItemsByOrderUID(ctx, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения items для заказа %s: %w", order.OrderUID, err)
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по заказам: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) getDeliveryByOrderUID(ctx context.Context, orderUID string) (model.Delivery, error) {
	var delivery model.Delivery
	query := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`
	err := r.db.QueryRow(ctx, query, orderUID).Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	return delivery, err
}

func (r *OrderRepository) getPaymentByOrderUID(ctx context.Context, orderUID string) (model.Payment, error) {
	var payment model.Payment
	query := `SELECT transaction_uid, currency, provider, amount, payment_dt, bank, delivery_cost FROM payment WHERE order_uid = $1`
	err := r.db.QueryRow(ctx, query, orderUID).Scan(&payment.Transaction, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost)
	return payment, err
}

func (r *OrderRepository) getItemsByOrderUID(ctx context.Context, orderUID string) ([]model.Item, error) {
	var items []model.Item
	query := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand FROM items WHERE order_uid = $1`
	rows, err := r.db.Query(ctx, query, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
