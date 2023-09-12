package usecase

import (
	"context"

	"github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type RestoreCacheUseCase struct {
	orderRepository domain.Repository
	cacheRepository domain.CacheRepository
}

func NewRestoreCacheUseCase(
	orderRepository domain.Repository,
	cacheRepository domain.CacheRepository,
) *RestoreCacheUseCase {
	return &RestoreCacheUseCase{
		orderRepository: orderRepository,
		cacheRepository: cacheRepository,
	}
}

func (uc *RestoreCacheUseCase) Execute(ctx context.Context) error {
	orders, err := uc.orderRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		_, err := uc.cacheRepository.SaveOrder(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}
