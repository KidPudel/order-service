package order

import (
	"github.com/samber/mo"

	pb "github.com/KidPudel/order-service/proto/order"
)

type OrderInfo struct {
	Type    uint32
	Amount  uint32
	Comment mo.Option[string]
}

func NewOrderInfo(pbOrderInfo *pb.OrderInfo) OrderInfo {
	comment := mo.None[string]()
	if pbOrderInfo.GetComment() != "" {
		comment = mo.Some(pbOrderInfo.GetComment())
	}

	return OrderInfo{
		Type:    pbOrderInfo.GetType(),
		Amount:  pbOrderInfo.GetAmount(),
		Comment: comment,
	}
}
