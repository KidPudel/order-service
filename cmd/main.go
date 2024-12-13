package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	handlers "github.com/KidPudel/order-service/internal/adapters/grpc"
	orderRepositories "github.com/KidPudel/order-service/internal/adapters/repositories/order"
	"github.com/KidPudel/order-service/internal/infrastructure/redis"
	orderUsecases "github.com/KidPudel/order-service/internal/usecases/order"
	pbDelivery "github.com/KidPudel/order-service/proto/delivery"
	pbOrder "github.com/KidPudel/order-service/proto/order"
)

func main() {
	listenConfig, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	// grpc server!
	server := grpc.NewServer()

	// grpc client
	clientConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	// client stub
	client := pbDelivery.NewDeliveryClient(clientConn)

	ctx, cancel := context.WithCancel(context.Background())

	// db
	redisClient := redis.NewRedis()

	// repositories
	orderRepository := orderRepositories.NewOrderRepostory(redisClient)

	// usecases
	orderUsecase := orderUsecases.NewOrderUsecase(ctx, orderUsecases.OrderUsecaseOptions{
		OrderRepository: orderRepository,
		DeliveryClient:  client,
	})

	// handlers
	opts := handlers.OrderOptions{
		OrderUsecase: orderUsecase,
	}

	pbOrder.RegisterOrderServer(server, handlers.NewOrderServer(opts))

	if err := server.Serve(listenConfig); err != nil {
		log.Fatal("failed to serve")
	}

	log.Println("end")

	cancel()

}
