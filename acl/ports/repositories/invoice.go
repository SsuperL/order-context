package repositories

import (
	"order-service/acl/adapters/pl"
	"order-service/domain/aggregate"
	ohs_pl "order-service/ohs/local/pl"
)

// InvoiceRepository 发票资源库端口，定义操作领域资源的方法,依赖倒置
type InvoiceRepository interface {
	GetInvoiceDetail(invoiceID, siteCode string) (pl.Invoice, error)
	GetInvoiceList(ohs_pl.ListInvoiceParams) ([]pl.Invoice, int, error)
	CreateInvoice(*aggregate.InvoiceAggregate, string) error
	UpdateInvoice(invoiceID, siteCode string, params ohs_pl.UpdateInvoiceParams) error
	CheckInvoiceExists(invoiceID, siteCode string) error
}
