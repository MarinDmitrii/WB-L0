package usecase

import (
	"context"

	"github.com/MarinDmitrii/WB-L0/internal/order/domain"
)

type GetOrderByIDUseCase struct {
	orderRepository domain.Repository
}

func NewGetOrderByIDUseCase(orderRepository domain.Repository) *GetOrderByIDUseCase {
	return &GetOrderByIDUseCase{orderRepository: orderRepository}
}

func (uc *GetOrderByIDUseCase) Execute(ctx context.Context, orderUID string) (domain.Order, error) {
	return uc.orderRepository.GetOrderByID(ctx, orderUID)
}
