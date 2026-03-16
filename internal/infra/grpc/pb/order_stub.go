package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
)

type UnimplementedOrderServiceServer struct{}

type CreateOrderRequest struct {
	Id    string
	Price float32
	Tax   float32
}

type CreateOrderResponse struct {
	Id         string
	Price      float32
	Tax        float32
	FinalPrice float32
}

type ListOrdersRequest struct{}

type Order struct {
	Id         string
	Price      float32
	Tax        float32
	FinalPrice float32
}

type ListOrdersResponse struct {
	Orders []*Order
}

type OrderServiceServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	ListOrders(context.Context, *ListOrdersRequest) (*ListOrdersResponse, error)
}

func RegisterOrderServiceServer(s *grpc.Server, srv OrderServiceServer) {}

func (UnimplementedOrderServiceServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, nil
}

func (UnimplementedOrderServiceServer) ListOrders(context.Context, *ListOrdersRequest) (*ListOrdersResponse, error) {
	return nil, nil
}
