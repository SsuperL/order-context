package repositories

import (
	"order-service/acl/adapters/pl"
	"order-service/common"
	"order-service/domain/aggregate"
)

// OrderRepository 订单资源库端口，定义操作领域资源的方法,依赖倒置
type OrderRepository interface {
	GetOrderDetail(orderID, siteCode string) (pl.Order, error)
	GetOrderList(common.ListOrderParams) ([]pl.Order, error)
	CreateOrder(*aggregate.AggregateRoot, string) error
	UpdateOrderStatus(orderID, siteCode string, status common.StatusType) error
	CheckOrderExists(orderID, siteCode string) (bool, error)
}
