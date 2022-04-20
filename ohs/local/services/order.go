package services

import (
	"fmt"
	dao "order-context/acl/adapters/repositories"
	"order-context/domain/aggregate"
	"order-context/domain/services"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
	"order-context/utils/common"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// CreateOrderAppService 创建订单
func CreateOrderAppService(id, siteCode string, options ...aggregate.RootOptions) (*pl.CreateOrderResponse, error) {
	orderAggregate := aggregate.NewOrderAggregateRoot(id, options...)
	orderService := services.NewOrderService(orderAggregate)
	orderID, err := orderService.CreateOrder(siteCode)
	if err != nil {
		return &pl.CreateOrderResponse{}, err
	}

	return &pl.CreateOrderResponse{
		Id: orderID,
	}, nil
}

// UpdateOrderAppService 更新订单
func UpdateOrderAppService(id, siteCode string, orderStatus common.StatusType, options ...aggregate.RootOptions) (*pl.UpdateOrderResponse, error) {
	orderAggregate := aggregate.NewOrderAggregateRoot(id, options...)
	orderService := services.NewOrderService(orderAggregate)
	err := orderService.UpdateOrderStatus(siteCode, orderStatus)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pl.UpdateOrderResponse{Success: false}, errors.OrderNotFound("Order not found")
		}

		return &pl.UpdateOrderResponse{Success: false}, status.Errorf(codes.Internal, "Error updating order status: %v", err)
	}
	return &pl.UpdateOrderResponse{Success: true}, nil
}

// GetOrderDetailAppService 获取订单详情
func GetOrderDetailAppService(id, siteCode string) (*pl.GetOrderDetailResponse, error) {
	orderDAO := dao.NewOrderAdapter()
	order, err := orderDAO.GetOrderDetail(id, siteCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pl.GetOrderDetailResponse{}, errors.OrderNotFound("Order not found")
		}
		return &pl.GetOrderDetailResponse{}, errors.InternalServerError(fmt.Sprintf("Get order detail error: %v", err))
	}
	return &pl.GetOrderDetailResponse{Result: pl.ToOrderBaseResponse(order)}, nil
}

// GetOrderListAppService 获取订单列表
func GetOrderListAppService(params pl.ListOrderParams) (*pl.GetOrderListResponse, error) {
	orderDAO := dao.NewOrderAdapter()
	orders, total, err := orderDAO.GetOrderList(params)
	if err != nil {
		return &pl.GetOrderListResponse{}, errors.InternalServerError(fmt.Sprintf("get order list error: %v", err))
	}

	datas := make([]*pl.OrderBase, 0)
	for _, order := range orders {
		datas = append(datas, pl.ToOrderBaseResponse(order))
	}

	return &pl.GetOrderListResponse{
		Data:  datas,
		Total: int32(total),
	}, nil
}

// OrderDAOAppService 提供给同一进程的聚合调用
func OrderDAOAppService(id, siteCode string) (*pl.GetOrderDetailResponse, error) {
	order, err := GetOrderDetailAppService(id, siteCode)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// OrderAppService 订单本地服务
// type OrderAppService struct {
// 	// 领域服务
// 	orderService *services.OrderService
// 	// 读操作，数据访问
// 	orderDAO repositories.OrderRepository
// }

// // NewOrderAppService 订单本地服务构造函数
// func NewOrderAppService(id string, options ...aggregate.RootOptions) OrderAppService {
// 	root := aggregate.NewOrderAggregateRoot(id, options...)

// 	return OrderAppService{
// 		orderService: services.NewOrderService(root),
// 		orderDAO:     dao.NewOrderAdapter(),
// 	}
// }

// // OrderDAOAppService 提供给同一进程的聚合调用
// func OrderDAOAppService(rootID, siteCode string) (*pl.GetOrderDetailResponse, error) {
// 	// 创建订单服务
// 	orderAppService := NewOrderAppService(rootID)
// 	// 订单数据模型
// 	order, err := orderAppService.GetOrderDetail(siteCode)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return order, nil
// }

// // CreateOrder 创建订单
// func (a *OrderAppService) CreateOrder(siteCode string) (*pl.CreateOrderResponse, error) {
// 	orderID, err := a.orderService.CreateOrder(siteCode)
// 	if err != nil {
// 		return &pl.CreateOrderResponse{}, err
// 	}

// 	return &pl.CreateOrderResponse{
// 		Id: orderID,
// 	}, nil
// }

// // UpdateOrder 更新订单
// func (a *OrderAppService) UpdateOrder(siteCode string, orderStatus common.StatusType) (*pl.UpdateOrderResponse, error) {
// 	err := a.orderService.UpdateOrderStatus(siteCode, orderStatus)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return &pl.UpdateOrderResponse{Success: false}, errors.OrderNotFound("Order not found")
// 		}

// 		return &pl.UpdateOrderResponse{Success: false}, status.Errorf(codes.Internal, "Error updating order status: %v", err)
// 	}
// 	return &pl.UpdateOrderResponse{Success: true}, nil
// }

// // GetOrderDetail 获取订单详情
// func (a *OrderAppService) GetOrderDetail(siteCode string) (*pl.GetOrderDetailResponse, error) {
// 	order, err := a.orderDAO.GetOrderDetail(a.orderService.Order.Order.ID, siteCode)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return &pl.GetOrderDetailResponse{}, errors.OrderNotFound("Order not found")
// 		}
// 		return &pl.GetOrderDetailResponse{}, errors.InternalServerError(fmt.Sprintf("Get order detail error: %v", err))
// 	}
// 	return &pl.GetOrderDetailResponse{Result: pl.ToOrderBaseResponse(order)}, nil
// }

// // GetOrderList 获取订单列表
// func (a *OrderAppService) GetOrderList(params pl.ListOrderParams) (*pl.GetOrderListResponse, error) {
// 	orders, total, err := a.orderDAO.GetOrderList(params)
// 	if err != nil {
// 		return &pl.GetOrderListResponse{}, errors.InternalServerError(fmt.Sprintf("get order list error: %v", err))
// 	}

// 	datas := make([]*pl.OrderBase, 0)
// 	for _, order := range orders {
// 		datas = append(datas, pl.ToOrderBaseResponse(order))
// 	}

// 	return &pl.GetOrderListResponse{
// 		Data:  datas,
// 		Total: int32(total),
// 	}, nil
// }
