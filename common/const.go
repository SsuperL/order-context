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
	// End ...
	End
)

const (
	// Invoiced 已开票
	Invoiced InvoiceStatusType = iota + 1
	// UnInvoiced 未开票
	UnInvoiced
	// InvoiceEnd ...
	InvoiceEnd
)

// UUIDurl url of uuid-service
var UUIDurl = "http://192.168.119.30:8181/api/v1"
