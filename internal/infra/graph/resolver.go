package graph

import "github.com/williamcardozo/go-clean-arch/internal/usecase"

type Resolver struct{
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}
