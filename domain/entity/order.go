package entity

import "order-service/common"

// Order 订单实体
type Order struct {
	// 订单ID
	ID string
	// 订单状态
	Status common.StatusType
	// 订单总价
	Price float32
}
