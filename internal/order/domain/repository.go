package domain

import (
	"context"
)

type Repository interface {
	SaveOrder(ctx context.Context, order Order) (string, error)
	GetAll(ctx context.Context) ([]Order, error)
	GetOrderByUID(ctx context.Context, orderUID string) (Order, error)
}

type CacheRepository interface {
	SaveOrder(ctx context.Context, order Order) (string, error)
	GetOrderByUID(ctx context.Context, orderUID string) (Order, error)
}
