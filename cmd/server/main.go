package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/williamcardozo/go-clean-arch/internal/infra/database"
	"github.com/williamcardozo/go-clean-arch/internal/infra/graph"
	"github.com/williamcardozo/go-clean-arch/internal/infra/grpc/pb"
	"github.com/williamcardozo/go-clean-arch/internal/infra/grpc/service"
	"github.com/williamcardozo/go-clean-arch/internal/infra/web"
	"github.com/williamcardozo/go-clean-arch/internal/infra/web/webserver"
	"github.com/williamcardozo/go-clean-arch/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/orders")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = database.RunMigrations(db, "./migrations")
	if err != nil {
		panic(err)
	}

	orderRepository := database.NewOrderRepository(db)
	createOrderUseCase := *usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := *usecase.NewListOrdersUseCase(orderRepository)

	webserver := webserver.NewWebServer(":8000")
	webOrderHandler := web.NewOrderHandler(orderRepository)
	webserver.AddHandler("/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			webOrderHandler.Create(w, r)
		} else {
			webOrderHandler.List(w, r)
		}
	})
	fmt.Println("Starting web server on port 8000")
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port 50051")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port 8080")
	http.ListenAndServe(":8080", nil)
}
