package resources

import (
	"context"
	"order-context/domain/aggregate"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
	"order-context/ohs/local/services"
	"order-context/utils/common"
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

	res, err := services.CreateInvoiceAppService(rootID, "", getSiteCode(ctx), invoiceOption, orderOption)
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

	invoiceOption := aggregate.WithInvoiceOption(status, req.GetPath(), "")
	params := pl.UpdateInvoiceParams{
		Status: status,
		Path:   req.GetPath(),
	}
	res, err := services.UpdateInvoiceAppService(id, getSiteCode(ctx), params, invoiceOption)
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

	res, err := services.GetInvoiceDetailAppService(id, getSiteCode(ctx))
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

	params := pl.ListInvoiceParams{
		OrderID: req.GetOrderId(),
		Limit:   int(req.GetLimit()),
		Offset:  int(req.GetOffset()),
	}
	res, err := services.GetInvoiceListAppService(params)
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
