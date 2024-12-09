package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	handlers "github.com/KidPudel/order-service/internal/adapters/grpc"
	orderUsecases "github.com/KidPudel/order-service/internal/usecases/order"
	pb "github.com/KidPudel/order-service/proto/order"
)

func main() {
	listenConfig, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	// grpc server!
	server := grpc.NewServer()

	ctx, cancel := context.WithCancel(context.Background())

	// usecases
	orderUsecase := orderUsecases.NewOrderUsecase(ctx, orderUsecases.OrderUsecaseOptions{
		OrderRepository: nil,
	})

	// handlers
	opts := handlers.OrderOptions{
		OrderUsecase: orderUsecase,
	}

	pb.RegisterOrderServer(server, handlers.NewOrderServer(opts))

	if err := server.Serve(listenConfig); err != nil {
		log.Fatal("failed to serve")
	}

	log.Println("end")

	cancel()

}
