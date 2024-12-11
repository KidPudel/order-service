package order

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"

	orderModel "github.com/KidPudel/order-service/internal/entities/order"
	pbDelivery "github.com/KidPudel/order-service/proto/delivery"
)

const (
	workersNumber       = 3
	maxConcurrentOrders = 9
)

type OrderRepository interface {
	AddOrder(ctx context.Context, orderInfo orderModel.OrderInfo) error
}

type OrderUsecaseOptions struct {
	OrderRepository OrderRepository
	DeliveryClient  pbDelivery.DeliveryClient
}

type OrderUsecase struct {
	options           OrderUsecaseOptions
	orderRequestQueue chan orderModel.OrderInfo
}

func NewOrderUsecase(ctx context.Context, options OrderUsecaseOptions) OrderUsecase {
	orderRequestQueue := make(chan orderModel.OrderInfo, maxConcurrentOrders)

	orderUsecase := OrderUsecase{
		options:           options,
		orderRequestQueue: orderRequestQueue,
	}

	// concurrency baby
	for i := 0; i < workersNumber; i++ {
		go orderUsecase.orderWorker(ctx)
	}

	return orderUsecase
}

func (u OrderUsecase) MakeOrder(ctx context.Context, orderInfo orderModel.OrderInfo) error {
	// add to the todo queue
	u.orderRequestQueue <- orderInfo

	return nil
}

func (u OrderUsecase) orderWorker(ctx context.Context) {
	for {
		select {
		case orderRequest := <-u.orderRequestQueue:
			fmt.Println(orderRequest)
			// read to redis
			err := u.options.OrderRepository.AddOrder(ctx, orderRequest)
			if err != nil {
				log.Println("failed to add to redis: ", err)
			}
			// send to delivery (shortcut for sake of mvp)
			comment, _ := orderRequest.Comment.Get()
			u.options.DeliveryClient.SendToDelivery(ctx, &pbDelivery.OrderInfo{
				Type:    proto.Uint32(orderRequest.Type),
				Amount:  proto.Uint32(orderRequest.Amount),
				Comment: proto.String(comment),
			})
		case <-ctx.Done():
			log.Println("done handling your stupid orders, we are quiting")
			break
		}
	}
}
