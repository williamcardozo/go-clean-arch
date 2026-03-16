package web

import (
	"encoding/json"
	"net/http"

	"github.com/williamcardozo/go-clean-arch/internal/entity"
	"github.com/williamcardozo/go-clean-arch/internal/usecase"
)

type OrderHandler struct {
	OrderRepository       entity.OrderRepositoryInterface
	CreateOrderUseCase    usecase.CreateOrderUseCase
	ListOrdersUseCase     usecase.ListOrdersUseCase
}

func NewOrderHandler(repository entity.OrderRepositoryInterface) *OrderHandler {
	return &OrderHandler{
		OrderRepository:    repository,
		CreateOrderUseCase: *usecase.NewCreateOrderUseCase(repository),
		ListOrdersUseCase:  *usecase.NewListOrdersUseCase(repository),
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrderUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
