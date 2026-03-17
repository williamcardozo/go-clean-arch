package usecase

import (
	"github.com/williamcardozo/go-clean-arch/internal/entity"
)

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	return dto, nil
}
