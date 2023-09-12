package adapters

import (
	"context"
	"fmt"

	domain "github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type CacheOrderRepository struct {
	cache         map[string]domain.Order
	cache2        map[int]string
	autoIncrement int
	maxSize       int
}

func NewCacheOrderRepository(maxSize int) *CacheOrderRepository {
	return &CacheOrderRepository{
		cache:         make(map[string]domain.Order, maxSize),
		cache2:        make(map[int]string, maxSize),
		autoIncrement: 1,
		maxSize:       maxSize,
	}
}

func (r *CacheOrderRepository) SaveOrder(ctx context.Context, domainOrder domain.Order) (string, error) {
	if r.autoIncrement == r.maxSize {
		r.autoIncrement = 1
	}
	if len(r.cache) == r.maxSize {
		delOrder := r.cache2[r.autoIncrement]
		delete(r.cache, delOrder)
	}

	r.cache[domainOrder.OrderUID] = domainOrder
	r.cache2[r.autoIncrement] = domainOrder.OrderUID
	r.autoIncrement++
	return domainOrder.OrderUID, nil
}

func (r *CacheOrderRepository) GetOrderByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	if order, ok := r.cache[orderUID]; !ok {
		return domain.Order{}, fmt.Errorf("Error: can't find order in cache")
	} else {
		return order, nil
	}
}
