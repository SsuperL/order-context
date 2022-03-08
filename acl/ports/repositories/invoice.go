package repositories

import (
	// "order-service/acl/adapters/repositories"
	"order-service/acl/adapters/pl"
	"order-service/common"
	"order-service/domain/aggregate"
)

// InvoiceRepository 发票资源库端口，定义操作领域资源的方法,依赖倒置
type InvoiceRepository interface {
	GetInvoiceDetail(invoiceID, siteCode string) (pl.Invoice, error)
	GetInvoiceList(common.ListInvoiceParams) ([]pl.Invoice, error)
	CreateInvoice(*aggregate.InvoiceAggregate, string) error
	UpdateInvoice(invoiceID, siteCode string, params common.UpdateInvoiceParams) error
	CheckInvoiceExists(invoiceID, siteCode string) error
}