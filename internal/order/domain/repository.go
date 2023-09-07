package domain

import (
	"context"
)

type Repository interface {
	SaveOrder(ctx context.Context, order Order) (string, error)
	GetOrderByID(ctx context.Context, orderUID string) (Order, error)
}
