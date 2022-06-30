package ohttp

import (
	"fmt"
	dao "order-context/acl/adapters/repositories"
	"order-context/domain/aggregate"
	"order-context/domain/services"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
)

// GetOrderListAppService 获取订单列表
func GetOrderListAppService(params pl.ListOrderParams) ([]pl.OrderResponse, int, error) {
	orderDAO := dao.NewOrderAdapter()
	orders, total, err := orderDAO.GetOrderList(params)
	if err != nil {
		return nil, total, errors.InternalServerError(fmt.Sprintf("get order list error: %v", err))
	}

	datas := make([]pl.OrderResponse, 0)
	for _, order := range orders {
		datas = append(datas, pl.OrderResponse{
			ID:             order.ID,
			Status:         pl.OrderStatus(order.Status),
			Number:         order.Number,
			SpaceID:        order.SpaceID,
			PayID:          order.PayID,
			Price:          order.Price,
			PackageVersion: order.PackageVersion,
			PackagePrice:   order.PackagePrice,
			// SiteCode:       order.SiteCode,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		})
	}

	return datas, total, nil

}

// CreateOrderAppService 创建订单
func CreateOrderAppService(id, siteCode string, options ...aggregate.RootOptions) (pl.CreateOrderResult, error) {
	orderAggregate := aggregate.NewOrderAggregateRoot(id, options...)
	orderService := services.NewOrderService(orderAggregate)
	orderID, err := orderService.CreateOrder(siteCode)
	if err != nil {
		return pl.CreateOrderResult{}, err
	}

	return pl.CreateOrderResult{
		ID: orderID,
	}, nil
}
