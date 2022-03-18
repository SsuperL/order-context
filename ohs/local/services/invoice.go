package services

import (
	"fmt"
	dao "order-service/acl/adapters/repositories"
	"order-service/acl/ports/repositories"
	"order-service/domain/aggregate"
	"order-service/domain/services"
	"order-service/ohs/local/pl"
	"order-service/ohs/local/pl/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// InvoiceAppService 发票本地服务
type InvoiceAppService struct {
	invoiceService *services.InvoiceService
	invoiceDAO     repositories.InvoiceRepository
}

// NewInvoiceAppService 发票本地服务构造函数
func NewInvoiceAppService(rootID, id string, options ...aggregate.InvoiceOptions) InvoiceAppService {
	agInvoice := aggregate.NewInvoiceAggregate(rootID, id, options...)

	return InvoiceAppService{
		invoiceService: services.NewInvoiceService(agInvoice),
		invoiceDAO:     dao.NewInvoiceAdapter(),
	}
}

// CreateInvoice 创建发票
func (a *InvoiceAppService) CreateInvoice(siteCode string) (*pl.CreateInvoiceResponse, error) {
	//{
	//创建服务
	//获取订单的详情
	//订单数据模型
	order, err := OrderDAOAppService(a.invoiceService.InvoiceAg.RootID, siteCode)
	if err != nil {
		return &pl.CreateInvoiceResponse{}, err
	}
	a.invoiceService.InvoiceAg.Order.Price = order.Result.Price

	invoiceID, err := a.invoiceService.CreateInvoice(siteCode)
	if err != nil {
		return &pl.CreateInvoiceResponse{}, status.Errorf(codes.Internal, "Error creating invoice: %v", err)
	}

	return &pl.CreateInvoiceResponse{
		Id: invoiceID,
	}, nil
}

// UpdateInvoice 更新发票
func (a *InvoiceAppService) UpdateInvoice(siteCode string, params pl.UpdateInvoiceParams) (*pl.UpdateInvoiceResponse, error) {
	err := a.invoiceService.UpdateInvoice(siteCode, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pl.UpdateInvoiceResponse{Success: false}, errors.InvoiceNotFound("Order not found")
		}

		return &pl.UpdateInvoiceResponse{Success: false}, status.Errorf(codes.Internal, "Error updating order status: %v", err)
	}
	return &pl.UpdateInvoiceResponse{Success: true}, nil
}

// GetInvoiceDetail 获取发票详情
func (a *InvoiceAppService) GetInvoiceDetail(siteCode string) (*pl.GetInvoiceDetailResponse, error) {
	invoice, err := a.invoiceDAO.GetInvoiceDetail(a.invoiceService.InvoiceAg.Invoice.ID, siteCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pl.GetInvoiceDetailResponse{}, errors.InvoiceNotFound("Invoice not found")
		}
		return &pl.GetInvoiceDetailResponse{}, errors.InternalServerError(fmt.Sprintf("Get Invoice detail error: %v", err))
	}
	return &pl.GetInvoiceDetailResponse{Result: pl.ToInvoiceBaseResponse(invoice)}, nil
}

// GetInvoiceList 获取发票列表
func (a *InvoiceAppService) GetInvoiceList(params pl.ListInvoiceParams) (*pl.GetInvoiceListResponse, error) {
	invoices, total, err := a.invoiceDAO.GetInvoiceList(params)
	if err != nil {
		return &pl.GetInvoiceListResponse{}, errors.InternalServerError(fmt.Sprintf("get invoice list error: %v", err))
	}

	datas := make([]*pl.InvoiceBase, 0)
	for _, invoice := range invoices {
		datas = append(datas, pl.ToInvoiceBaseResponse(invoice))
	}

	return &pl.GetInvoiceListResponse{
		Data:  datas,
		Total: int32(total),
	}, nil
}
