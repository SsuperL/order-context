package services

import (
	"fmt"
	client_adapter "order-context/acl/adapters/clients"
	repository_adapter "order-context/acl/adapters/repositories"
	client_port "order-context/acl/ports/clients"
	repository_port "order-context/acl/ports/repositories"
	"order-context/domain/aggregate"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
)

// InvoiceService 发票领域服务
type InvoiceService struct {
	InvoicePort repository_port.InvoiceRepository
	OrderPort   repository_port.OrderRepository
	UUIDClient  client_port.UUIDClient
	InvoiceAg   *aggregate.InvoiceAggregate
}

// NewInvoiceService 发票领域服务构造函数
func NewInvoiceService(invoiceAg *aggregate.InvoiceAggregate) *InvoiceService {
	return &InvoiceService{
		InvoicePort: repository_adapter.NewInvoiceAdapter(),
		OrderPort:   repository_adapter.NewOrderAdapter(),
		UUIDClient:  client_adapter.NewUUIDAdapter(),
		InvoiceAg:   invoiceAg,
	}
}

// WithPortAndClientParams ...
func WithPortAndClientParams(invoiceAg *aggregate.InvoiceAggregate, orderPort repository_port.OrderRepository,
	invoicePort repository_port.InvoiceRepository, client client_port.UUIDClient) *InvoiceService {
	return &InvoiceService{
		InvoicePort: invoicePort,
		OrderPort:   orderPort,
		UUIDClient:  client,
		InvoiceAg:   invoiceAg,
	}
}

// CreateInvoice 创建发票
func (isv *InvoiceService) CreateInvoice(siteCode string) (invoiceID string, err error) {
	// 校验订单有效性
	exists, err := isv.OrderPort.CheckOrderExists(isv.InvoiceAg.RootID, siteCode)
	if err != nil {
		return
	} else if !exists {
		err = errors.OrderNotFound("Order not exists")
		return
	}
	// 生成发票id
	res, err := isv.UUIDClient.GetUUID(1)
	if err != nil {
		err = errors.InternalServerError(fmt.Sprintf("Get uuid failed: %v", err))
		return
	}
	invoiceID = res.ID
	// 创建发票
	isv.InvoiceAg.SetID(invoiceID)
	err = isv.InvoicePort.CreateInvoice(isv.InvoiceAg, siteCode)
	if err != nil {
		err = errors.InternalServerError(fmt.Sprintf("Create invoice failed: %v", err))
		return
	}

	return
}

// UpdateInvoice 更新发票
func (isv *InvoiceService) UpdateInvoice(siteCode string, params pl.UpdateInvoiceParams) error {
	// 校验发票有效性
	if err := isv.InvoicePort.CheckInvoiceExists(isv.InvoiceAg.Invoice.ID, siteCode); err != nil {
		return err
	}
	// 更新发票
	if err := isv.InvoicePort.UpdateInvoice(isv.InvoiceAg.Invoice.ID, siteCode, params); err != nil {
		return err
	}
	return nil
}

// GetInvoiceDetail 获取发票详情
// func (isv *InvoiceService) GetInvoiceDetail(siteCode string) (model.Invoice, error) {
// 	// 获取发票详情
// 	invoice, err := isv.InvoicePort.GetInvoiceDetail(isv.InvoiceAg.Invoice.ID, siteCode)
// 	if err != nil {
// 		return model.Invoice{}, err
// 	}
// 	return invoice, nil
// }

// // GetInvoiceList 获取发票列表
// func (isv *InvoiceService) GetInvoiceList(args ohs_pl.ListInvoiceParams) ([]model.Invoice, int, error) {
// 	// 获取发票列表
// 	invoices, total, err := isv.InvoicePort.GetInvoiceList(args)
// 	if err != nil {
// 		return nil, total, err
// 	}
// 	return invoices, total, nil
// }
