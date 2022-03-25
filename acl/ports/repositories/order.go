package repositories

import (
	"order-context/acl/adapters/pl"
	"order-context/common"
	"order-context/domain/aggregate"
	ohs_pl "order-context/ohs/local/pl"
)

// OrderRepository 订单资源库端口，定义操作领域资源的方法,依赖倒置
type OrderRepository interface {
	GetOrderDetail(orderID, siteCode string) (pl.Order, error)
	GetOrderList(ohs_pl.ListOrderParams) ([]pl.Order, int, error)
	CreateOrder(*aggregate.AggregateRoot, string) error
	UpdateOrderStatus(orderID, siteCode string, status common.StatusType) error
	CheckOrderExists(orderID, siteCode string) (bool, error)
}
