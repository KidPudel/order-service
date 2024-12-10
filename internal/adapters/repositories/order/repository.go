package order

import (
	"context"
	"fmt"

	orderModel "github.com/KidPudel/order-service/internal/entities/order"
	db "github.com/KidPudel/order-service/internal/infrastructure/redis"
)

type OrderRepository struct {
	db db.RedisClient
}

func NewOrderRepostory(db db.RedisClient) OrderRepository {
	return OrderRepository{
		db: db,
	}
}

func (repository OrderRepository) AddOrder(ctx context.Context, orderInfo orderModel.OrderInfo) error {
	id, err := repository.db.Client.Incr(ctx, "order:id:counter").Result()
	if err != nil {
		return err
	}
	comment, _ := orderInfo.Comment.Get()
	err = repository.db.Client.HSet(ctx, fmt.Sprintf("order:%d", id), map[string]string{"type": string(orderInfo.Type), "amount": string(orderInfo.Amount), "comment": comment}).Err()
	if err != nil {
		return err
	}
	return nil
}
