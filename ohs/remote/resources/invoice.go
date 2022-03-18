package resources

import (
	"context"
	"order-service/common"
	"order-service/domain/aggregate"
	"order-service/ohs/local/pl"
	"order-service/ohs/local/pl/errors"
	"order-service/ohs/local/services"
)

// InvoiceResource ...
type InvoiceResource struct {
	pl.UnimplementedInvoiceServiceServer
}

// NewInvoiceResource create instance of invoice resource
func NewInvoiceResource() pl.InvoiceServiceServer {
	return &InvoiceResource{}
}

// CreateInvoice handler request of invoice
func (i *InvoiceResource) CreateInvoice(ctx context.Context, req *pl.CreateInvoiceRequest) (*pl.CreateInvoiceResponse, error) {
	rootID := req.GetOrderId()
	status := common.InvoiceStatusType(req.GetStatus())
	if err := validateInvoiceStatus(status); err != nil {
		return nil, err
	}

	path := req.GetPath()
	name := req.GetName()
	price := req.GetPrice()
	invoiceOption := aggregate.WithInvoiceOption(status, path, name)
	orderOption := aggregate.WithOrderOptionForInvoice(price)

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	invoiceAppService := services.NewInvoiceAppService(rootID, "", invoiceOption, orderOption)

	res, err := invoiceAppService.CreateInvoice(getSiteCode(ctx))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateInvoice handle request of update invoice
func (i *InvoiceResource) UpdateInvoice(ctx context.Context, req *pl.UpdateInvoiceRequest) (*pl.UpdateInvoiceResponse, error) {
	id := req.GetId()
	status := common.InvoiceStatusType(req.GetStatus())
	if err := validateInvoiceStatus(status); err != nil {
		return nil, err
	}

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	invoiceAppService := services.NewInvoiceAppService("", id)
	params := pl.UpdateInvoiceParams{
		Status: status,
		Path:   req.GetPath(),
	}
	res, err := invoiceAppService.UpdateInvoice(getSiteCode(ctx), params)
	if err != nil {
		return nil, err
	}

	return res, nil

}

// GetInvoiceDetail handle request of get invoice detail
func (i *InvoiceResource) GetInvoiceDetail(ctx context.Context, req *pl.GetInvoiceDetailRequest) (*pl.GetInvoiceDetailResponse, error) {
	id := req.GetId()

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	invoiceAppService := services.NewInvoiceAppService("", id)

	res, err := invoiceAppService.GetInvoiceDetail(getSiteCode(ctx))
	if err != nil {
		return nil, err
	}

	return res, err
}

// GetInvoiceList handle request of get invoice list
func (i *InvoiceResource) GetInvoiceList(ctx context.Context, req *pl.GetInvoiceListRequest) (*pl.GetInvoiceListResponse, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	invoiceAppService := services.NewInvoiceAppService("", "")

	params := pl.ListInvoiceParams{
		OrderID: req.GetOrderId(),
		Limit:   int(req.GetLimit()),
		Offset:  int(req.GetOffset()),
	}
	res, err := invoiceAppService.GetInvoiceList(params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func validateInvoiceStatus(status common.InvoiceStatusType) error {
	if status >= common.InvoiceEnd {
		return errors.BadRequest("Invalid status")
	}
	return nil
}
