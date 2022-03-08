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

// NewInvoiceAggregate 发票聚合构造函数
func NewInvoiceAggregate(RootID, id, name, code, path string, status common.InvoiceStatusType, price float64) *InvoiceAggregate {
	return &InvoiceAggregate{
		RootID: id,
		Invoice: &entity.Invoice{
			ID:     id,
			Status: status,
			Path:   path,
			Detail: vo.InvoiceDetail{
				Name: name,
				Code: code,
			},
		},
		Order: vo.Order{
			Price: price,
		},
	}
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
