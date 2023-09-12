package usecase

import (
	"context"

	domain "github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type SaveOrderUseCase struct {
	orderRepository domain.Repository
	cacheRepository domain.CacheRepository
}

func NewSaveOrderUseCase(
	orderRepository domain.Repository,
	cacheRepository domain.CacheRepository,
) *SaveOrderUseCase {
	return &SaveOrderUseCase{
		orderRepository: orderRepository,
		cacheRepository: cacheRepository,
	}
}

func (uc *SaveOrderUseCase) Execute(ctx context.Context, order domain.Order) (string, error) {
	orderUID, err := uc.orderRepository.SaveOrder(ctx, order)
	if err != nil {
		return "", err
	}

	_, err = uc.cacheRepository.SaveOrder(ctx, order)
	if err != nil {
		return "", err
	}

	return orderUID, nil
}
