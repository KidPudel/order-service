package grpc

import (
	"context"
	"log"

	"google.golang.org/protobuf/proto"

	orderModel "github.com/KidPudel/order-service/internal/entities/order"
	pb "github.com/KidPudel/order-service/proto/order"
)

type OrderUsecase interface {
	MakeOrder(ctx context.Context, orderInfo orderModel.OrderInfo) error
}

type OrderOptions struct {
	OrderUsecase OrderUsecase
}

type OrderServer struct {
	pb.OrderServer
	options OrderOptions
}

func NewOrderServer(options OrderOptions) *OrderServer {
	return &OrderServer{
		options: options,
	}
}

func (orderService *OrderServer) MakeOrder(ctx context.Context, orderInfo *pb.OrderInfo) (*pb.OrderAcknowledgment, error) {
	log.Println("got order")

	order := orderModel.NewOrderInfo(orderInfo)

	err := orderService.options.OrderUsecase.MakeOrder(ctx, order)
	if err != nil {
		return &pb.OrderAcknowledgment{
			Response: proto.String(err.Error()),
		}, err
	}

	return &pb.OrderAcknowledgment{
		Response: proto.String("order created"),
	}, nil
}
