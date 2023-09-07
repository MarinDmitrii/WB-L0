package usecase

import (
	"context"

	domain "github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type SaveOrder struct {
	Order domain.Order
}

type SaveOrderUseCase struct {
	orderRepository domain.Repository
}

func NewSaveOrderUseCase(
	orderRepository domain.Repository,
) *SaveOrderUseCase {
	return &SaveOrderUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *SaveOrderUseCase) Execute(ctx context.Context, saveOrder *SaveOrder) (string, error) {
	orderUID, err := uc.orderRepository.SaveOrder(ctx, *&saveOrder.Order)
	if err != nil {
		return "", err
	}

	return orderUID, nil
}
