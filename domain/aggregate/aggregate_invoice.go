package aggregate

import (
	"order-service/common"
	"order-service/domain/entity"
	"order-service/domain/vo"
)

// InvoiceAggregate 发票聚合
type InvoiceAggregate struct {
	// 聚合根ID，全局唯一
	RootID string
	// 发票实体
	Invoice *entity.Invoice
	// 订单值对象
	Order vo.Order
}

// InvoiceOptions ...
type InvoiceOptions func(*InvoiceAggregate)

// WithInvoiceOption init invoice aggregate with invoice params
func WithInvoiceOption(status common.InvoiceStatusType, path, name string) InvoiceOptions {
	return func(i *InvoiceAggregate) {
		i.Invoice.Status = status
		i.Invoice.Path = path
		i.Invoice.Detail.Name = name
	}
}

// WithOrderOptionForInvoice init invoice aggregate with order params
func WithOrderOptionForInvoice(price float32) InvoiceOptions {
	return func(i *InvoiceAggregate) {
		i.Order.Price = price
	}
}

// NewInvoiceAggregate 发票聚合构造函数
func NewInvoiceAggregate(rootID, id string, options ...InvoiceOptions) *InvoiceAggregate {
	invoiceAg := &InvoiceAggregate{
		RootID:  rootID,
		Invoice: &entity.Invoice{},
	}
	for _, option := range options {
		option(invoiceAg)
	}
	return invoiceAg
}

// SetID 设置发票id
func (ia *InvoiceAggregate) SetID(id string) {
	ia.Invoice.ID = id
}

// GetID 获取发票id
func (ia *InvoiceAggregate) GetID() string {
	return ia.Invoice.ID
}

// GetStatus 获取发票状态
func (ia *InvoiceAggregate) GetStatus() common.InvoiceStatusType {
	return ia.Invoice.Status
}

// SetStatus 设置发票状态
func (ia *InvoiceAggregate) SetStatus(status common.InvoiceStatusType) {
	ia.Invoice.Status = status
}

// SetPath 设置发票保存路径
func (ia *InvoiceAggregate) SetPath(path string) {
	ia.Invoice.Path = path
}

// GetPath 获取发票保存路径
func (ia *InvoiceAggregate) GetPath() string {
	return ia.Invoice.Path
}
