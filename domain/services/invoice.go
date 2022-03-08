package services

import (
	client_adapter "order-service/acl/adapters/clients"
	model "order-service/acl/adapters/pl"
	repository_adapter "order-service/acl/adapters/repositories"
	client_port "order-service/acl/ports/clients"
	repository_port "order-service/acl/ports/repositories"
	"order-service/common"
	"order-service/domain/aggregate"

	"gorm.io/gorm"
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

// CreateInvoice 创建发票
func (isv *InvoiceService) CreateInvoice(siteCode string) (invoiceID string, err error) {
	// 校验订单有效性
	exists, err := isv.OrderPort.CheckOrderExists(isv.InvoiceAg.RootID, siteCode)
	if err != nil {
		return
	} else if !exists {
		// TODO: 错误处理
		return
	}
	// 生成发票id
	res, err := isv.UUIDClient.GetUUID(1)
	invoiceID = res.ID
	if err != nil {
		return
	}
	// 创建发票
	isv.InvoiceAg.SetID(invoiceID)
	err = isv.InvoicePort.CreateInvoice(isv.InvoiceAg, siteCode)
	if err != nil {
		return
	}

	return
}

// UpdateInvoice 更新发票
func (isv *InvoiceService) UpdateInvoice(siteCode string) error {
	// 校验发票有效性
	if err := isv.InvoicePort.CheckInvoiceExists(isv.InvoiceAg.Invoice.ID, siteCode); err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}
	param := common.UpdateInvoiceParams{
		Status: isv.InvoiceAg.GetStatus(),
		Path:   isv.InvoiceAg.GetPath(),
	}
	// 更新发票
	if err := isv.InvoicePort.UpdateInvoice(isv.InvoiceAg.Invoice.ID, siteCode, param); err != nil {
		return err
	}
	return nil
}

// GetInvoiceDetail 获取发票详情
func (isv *InvoiceService) GetInvoiceDetail(siteCode string) (model.Invoice, error) {
	// 获取发票详情
	invoice, err := isv.InvoicePort.GetInvoiceDetail(isv.InvoiceAg.Invoice.ID, siteCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Invoice{}, err
		}
		return model.Invoice{}, err
	}
	return invoice, nil
}

// GetInvoiceList 获取发票列表
func (isv *InvoiceService) GetInvoiceList(args common.ListInvoiceParams) ([]model.Invoice, error) {
	// 获取发票列表
	invoices, err := isv.InvoicePort.GetInvoiceList(args)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}
