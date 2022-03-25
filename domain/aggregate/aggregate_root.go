package aggregate

import (
	"order-context/common"
	"order-context/domain/entity"
	"order-context/domain/vo"
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

// RootOptions ...
type RootOptions func(*AggregateRoot)

// WithOrderOption Init Order with order params
func WithOrderOption(status common.StatusType, price float32) RootOptions {
	return func(ag *AggregateRoot) {
		ag.Order.Status = status
		ag.Order.Price = price
	}
}

// WithPaymentOption init order with payment params
func WithPaymentOption(voucher, source, currency string, total float32) RootOptions {
	return func(ag *AggregateRoot) {
		ag.Payment.Voucher = voucher
		ag.Payment.Source = source
		ag.Payment.Currency = currency
		ag.Payment.Total = total
	}
}

// WithPackageOption init order with package params
func WithPackageOption(version string, price float32) RootOptions {
	return func(ag *AggregateRoot) {
		ag.Package.Version = version
		ag.Package.Price = price
	}
}

// WithSpaceOption init order with space id
func WithSpaceOption(spaceID string) RootOptions {
	return func(ag *AggregateRoot) {
		ag.Space.ID = spaceID
	}
}

// NewOrderAggregateRoot 聚合根构造函数，在本地应用服务生成
func NewOrderAggregateRoot(id string, options ...RootOptions) *AggregateRoot {
	root := &AggregateRoot{
		Order: &entity.Order{ID: id},
		Space: &entity.Space{},
	}
	for _, op := range options {
		op(root)
	}
	return root
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
