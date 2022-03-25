package pl

import (
	model "order-context/acl/adapters/pl"
	"order-context/common"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListOrderParams 获取订单列表参数，待放置pl
type ListOrderParams struct {
	SpaceID string
	Status  common.StatusType
	Limit   int
	Offset  int
}

// UpdateInvoiceParams 更新发票参数，带放置pl
type UpdateInvoiceParams struct {
	Status common.InvoiceStatusType
	Path   string
}

// ListInvoiceParams 获取发票列表参数，待放置pl
type ListInvoiceParams struct {
	OrderID string
	Status  common.InvoiceStatusType
	Limit   int
	Offset  int
}

// ToOrderBaseResponse 将数据库模型对象转换为消息契约模型对象
func ToOrderBaseResponse(order model.Order) *OrderBase {
	result := &OrderBase{
		Id:             order.ID,
		Status:         OrderStatus(order.Status),
		Number:         order.Number,
		SpaceId:        order.SpaceID,
		PayId:          order.PayID,
		Price:          order.Price,
		PackageVersion: order.PackageVersion,
		PackagePrice:   order.PackagePrice,
		SiteCode:       order.SiteCode,
		CreatedAt:      timestamppb.New(order.CreatedAt),
		UpdatedAt:      timestamppb.New(order.UpdatedAt),
	}
	return result
}

// ToInvoiceBaseResponse 将数据库模型对象转换为消息契约模型对象
func ToInvoiceBaseResponse(invoice model.Invoice) *InvoiceBase {
	result := &InvoiceBase{
		Id:        invoice.ID,
		OrderId:   invoice.OrderID,
		Price:     invoice.Price,
		Name:      invoice.Name,
		Code:      invoice.Code,
		Status:    InvoiceStatus(invoice.Status),
		Path:      invoice.Path,
		SiteCode:  invoice.SiteCode,
		CreatedAt: timestamppb.New(invoice.CreatedAt),
		UpdatedAt: timestamppb.New(invoice.UpdatedAt),
	}
	return result
}
