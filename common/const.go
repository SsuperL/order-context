package common

// StatusType 订单支付状态
type StatusType int32

// InvoiceStatusType 发票支付状态
type InvoiceStatusType int32

const (
	// Unpaid 未支付
	Unpaid StatusType = iota + 1
	// Paid 已支付
	Paid
	// Canceled 已取消
	Canceled
	// Failed 已失效
	Failed
	// Invoiced 已开票
	Invoiced InvoiceStatusType = iota + 1
	// UnInvoiced 未开票
	UnInvoiced
)

// ListOrderParams 获取订单列表参数，待放置pl
type ListOrderParams struct {
	SpaceID string
	Status  StatusType
	Limit   int
	Offset  int
}

// UpdateInvoiceParams 更新发票参数，带放置pl
type UpdateInvoiceParams struct {
	Status InvoiceStatusType
	Path   string
}

// ListInvoiceParams 获取发票列表参数，待放置pl
type ListInvoiceParams struct {
	OrderID string
	Status  InvoiceStatusType
	Limit   int
	Offset  int
}

// UUIDurl url of uuid-service
var UUIDurl = "http://127.0.0.1:8181/api/v1"
