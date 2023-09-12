package usecase

import (
	"context"

	"github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type GetOrderByUIDUseCase struct {
	orderRepository domain.Repository
	cacheRepository domain.CacheRepository
}

func NewGetOrderByUIDUseCase(
	orderRepository domain.Repository,
	cacheRepository domain.CacheRepository,
) *GetOrderByUIDUseCase {
	return &GetOrderByUIDUseCase{
		orderRepository: orderRepository,
		cacheRepository: cacheRepository,
	}
}

func (uc *GetOrderByUIDUseCase) Execute(ctx context.Context, orderUID string) (domain.Order, error) {
	order, err := uc.cacheRepository.GetOrderByUID(ctx, orderUID)
	if err != nil {
		order, err = uc.orderRepository.GetOrderByUID(ctx, orderUID)
		if err != nil {
			return domain.Order{}, err
		}
		uc.cacheRepository.SaveOrder(ctx, order)
	}

	return order, nil
}
