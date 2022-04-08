package resources

import (
	"context"
	"order-context/domain/aggregate"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
	"order-context/ohs/local/services"
	"order-context/utils/common"

	"google.golang.org/grpc/metadata"
)

// OrderResource 订单资源
type OrderResource struct {
	pl.UnimplementedOrderServiceServer
}

// NewOrderResource create instance of order resource
func NewOrderResource() pl.OrderServiceServer {
	return &OrderResource{}
}

// CreateOrder handle request of create order
func (r *OrderResource) CreateOrder(ctx context.Context, req *pl.CreateOrderRequest) (*pl.CreateOrderResponse, error) {
	status := common.StatusType(req.GetStatus())
	if err := validateStatus(status); err != nil {
		return nil, err
	}
	price := req.GetPrice()
	packageVersion := req.GetPackageVersion()
	packagePrice := req.GetPackagePrice()
	spaceID := req.GetSpaceId()

	orderOption := aggregate.WithOrderOption(status, price)
	spaceOption := aggregate.WithSpaceOption(spaceID)
	packageOption := aggregate.WithPackageOption(packageVersion, packagePrice)

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	orderAppService := services.NewOrderAppService("", orderOption, spaceOption, packageOption)
	response, err := orderAppService.CreateOrder(getSiteCode(ctx))
	if err != nil {
		return &pl.CreateOrderResponse{}, err
	}
	return response, nil

}

// UpdateOrder update order
func (r *OrderResource) UpdateOrder(ctx context.Context, req *pl.UpdateOrderRequest) (*pl.UpdateOrderResponse, error) {
	id := req.GetId()
	status := common.StatusType(req.GetStatus())
	if err := validateStatus(status); err != nil {
		return nil, err
	}

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	orderAppService := services.NewOrderAppService(id)

	res, err := orderAppService.UpdateOrder(getSiteCode(ctx), status)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetOrderDetail get order details
func (r *OrderResource) GetOrderDetail(ctx context.Context, req *pl.GetOrderDetailRequest) (*pl.GetOrderDetailResponse, error) {
	id := req.GetId()

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	orderAppService := services.NewOrderAppService(id)

	res, err := orderAppService.GetOrderDetail(getSiteCode(ctx))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetOrderList get order list
func (r *OrderResource) GetOrderList(ctx context.Context, req *pl.GetOrderListRequest) (*pl.GetOrderListResponse, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	orderAppService := services.NewOrderAppService("")

	status := common.StatusType(req.GetStatus())
	if err := validateStatus(status); err != nil {
		return nil, err
	}

	params := pl.ListOrderParams{
		SpaceID: req.GetSpaceId(),
		Status:  status,
		Limit:   int(req.GetLimit()),
		Offset:  int(req.GetOffset()),
	}

	res, err := orderAppService.GetOrderList(params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func validateStatus(status common.StatusType) error {
	if status >= common.End {
		return errors.BadRequest("Invalid status")
	}
	return nil
}

func getSiteCode(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	return md["site-code"][0]
}

// 客户端断开请求及超时处理
func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return errors.Canceled("request is canceled")
	case context.DeadlineExceeded:
		return errors.DeadlineExceeded("deadline is exceeded")
	default:
		return nil
	}
}
