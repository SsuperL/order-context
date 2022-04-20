package services

import (
	"fmt"
	dao "order-context/acl/adapters/repositories"
	"order-context/domain/aggregate"
	"order-context/domain/services"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// CreateInvoiceAppService 创建发票
func CreateInvoiceAppService(rootID, id, siteCode string, options ...aggregate.InvoiceOptions) (*pl.CreateInvoiceResponse, error) {
	invoiceAggregate := aggregate.NewInvoiceAggregate(rootID, id, options...)
	invoiceService := services.NewInvoiceService(invoiceAggregate)
	order, err := OrderDAOAppService(rootID, siteCode)
	if err != nil {
		return &pl.CreateInvoiceResponse{}, err
	}

	invoiceService.InvoiceAg.Order.Price = order.Result.Price

	invoiceID, err := invoiceService.CreateInvoice(siteCode)
	if err != nil {
		return &pl.CreateInvoiceResponse{}, err
	}

	return &pl.CreateInvoiceResponse{
		Id: invoiceID,
	}, nil
}

// UpdateInvoiceAppService 更新发票
func UpdateInvoiceAppService(id, siteCode string, params pl.UpdateInvoiceParams, options ...aggregate.InvoiceOptions) (*pl.UpdateInvoiceResponse, error) {
	invoiceAggregate := aggregate.NewInvoiceAggregate("", id, options...)
	invoiceService := services.NewInvoiceService(invoiceAggregate)
	err := invoiceService.UpdateInvoice(siteCode, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pl.UpdateInvoiceResponse{Success: false}, errors.InvoiceNotFound("Order not found")
		}

		return &pl.UpdateInvoiceResponse{Success: false}, status.Errorf(codes.Internal, "Error updating order status: %v", err)
	}
	return &pl.UpdateInvoiceResponse{Success: true}, nil
}

// GetInvoiceDetailAppService 获取发票详情
func GetInvoiceDetailAppService(id, siteCode string, options ...aggregate.InvoiceOptions) (*pl.GetInvoiceDetailResponse, error) {
	invoiceDAO := dao.NewInvoiceAdapter()
	invoice, err := invoiceDAO.GetInvoiceDetail(id, siteCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pl.GetInvoiceDetailResponse{}, errors.InvoiceNotFound("Invoice not found")
		}
		return &pl.GetInvoiceDetailResponse{}, errors.InternalServerError(fmt.Sprintf("Get Invoice detail error: %v", err))
	}
	return &pl.GetInvoiceDetailResponse{Result: pl.ToInvoiceBaseResponse(invoice)}, nil
}

// GetInvoiceList 获取发票列表
func GetInvoiceListAppService(params pl.ListInvoiceParams) (*pl.GetInvoiceListResponse, error) {
	invoiceDAO := dao.NewInvoiceAdapter()
	invoices, total, err := invoiceDAO.GetInvoiceList(params)
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
