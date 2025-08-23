package cache

import (
	"github.com/Booba186/level0/internal/model"
	"sync"
)

type Cache struct {
	mx     sync.RWMutex
	orders map[string]model.Order
}

func New() *Cache {
	return &Cache{
		orders: make(map[string]model.Order),
	}
}

func (c *Cache) Set(order model.Order) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) Get(orderUID string) (model.Order, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	order, found := c.orders[orderUID]
	return order, found
}
