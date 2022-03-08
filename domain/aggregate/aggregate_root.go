package aggregate

import (
	"order-service/common"
	"order-service/domain/entity"
	"order-service/domain/vo"
)

// AggregateRoot 订单聚合，聚合根
type AggregateRoot struct {
	// Order 订单实体
	Order *entity.Order
	// Payment 支付值对象
	Payment vo.Payment
	// Package 套餐值对象
	Package vo.Package
	// space 空间实体
	Space *entity.Space
}

// NewOrderAggregateRoot 聚合根构造函数，在本地应用服务生成
func NewOrderAggregateRoot(id, spaceID, voucher, source, version, currency string,
	status common.StatusType, price, total float64) *AggregateRoot {
	return &AggregateRoot{
		Order: &entity.Order{
			ID:     id,
			Status: status,
			Price:  total,
		},
		Payment: vo.Payment{
			Voucher:  voucher,
			Source:   source,
			Currency: currency,
			Total:    total,
		},
		Package: vo.Package{
			Version: version,
			Price:   price,
		},
		Space: &entity.Space{
			ID: spaceID,
		}}
}

// SetID 设置聚合根id
func (oar *AggregateRoot) SetID(id string) {
	oar.Order.ID = id
}

// GetID 获取聚合根id
func (oar *AggregateRoot) GetID() string {
	return oar.Order.ID
}

// GetStatus 获取订单状态
func (oar *AggregateRoot) GetStatus() common.StatusType {
	return oar.Order.Status
}

// SetStatus 设置订单状态
func (oar *AggregateRoot) SetStatus(status common.StatusType) {
	oar.Order.Status = status
}
