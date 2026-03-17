package graph

import (
	"context"

	"github.com/williamcardozo/go-clean-arch/internal/entity"
	"github.com/williamcardozo/go-clean-arch/internal/infra/graph/model"
	"github.com/williamcardozo/go-clean-arch/internal/usecase"
)

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderInput) (*entity.Order, error) {
	dto := usecase.OrderInputDTO{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}

	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &entity.Order{
		ID:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (r *queryResolver) Orders(ctx context.Context) ([]*entity.Order, error) {
	output, err := r.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var orders []*entity.Order
	for _, dto := range output {
		orders = append(orders, &entity.Order{
			ID:         dto.ID,
			Price:      dto.Price,
			Tax:        dto.Tax,
			FinalPrice: dto.FinalPrice,
		})
	}

	return orders, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
