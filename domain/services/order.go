package services

import (
	"fmt"
	client_adapter "order-context/acl/adapters/clients"
	repository_adapter "order-context/acl/adapters/repositories"
	client_port "order-context/acl/ports/clients"
	repository_port "order-context/acl/ports/repositories"
	"order-context/domain/aggregate"
	ohs_pl "order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
	"order-context/utils/common"
	"time"
)

// 后期修改为读取配置表
const expiration = 24

// OrderService 订单领域服务,对接聚合根，进行领域动作，可持有多个repository，接收依赖注入
type OrderService struct {
	// 资源库端口
	Port repository_port.OrderRepository
	// Log        common.Logger
	// 订单聚合根
	Order *aggregate.AggregateRoot
	// uuid服务客户端
	UUIDClient client_port.UUIDClient
}

// NewOrderService 构造函数，提供给应用服务
func NewOrderService(root *aggregate.AggregateRoot) *OrderService {
	return &OrderService{
		Port:       repository_adapter.NewOrderAdapter(),
		Order:      root,
		UUIDClient: client_adapter.NewUUIDAdapter(),
	}
}

// WithPortAndClient 传参方式构造函数，可替换
func WithPortAndClient(root *aggregate.AggregateRoot, port repository_port.OrderRepository, client client_port.UUIDClient) *OrderService {
	return &OrderService{
		Port:       port,
		Order:      root,
		UUIDClient: client,
	}
}

// CreateOrder 创建订单
func (osv *OrderService) CreateOrder(siteCode string) (orderID string, err error) {
	// 校验订单有效性
	args := ohs_pl.ListOrderParams{
		SpaceID: osv.Order.Space.ID,
		Status:  common.Unpaid,
		Limit:   1,
		Offset:  0,
	}
	orders, _, err := osv.Port.GetOrderList(args)
	if err != nil {
		return
	}

	// 校验是否存在未支付订单
	if len(orders) != 0 {
		// 订单是否过期判断
		order := orders[0]
		now := time.Now()
		if now.Sub(order.CreatedAt) >= time.Duration(expiration*time.Hour) {
			// 存在未支付订单已过期，更新订单状态
			if err = osv.UpdateOrderStatus(siteCode, common.Failed); err != nil {
				return
			}
		} else {
			// 未过期提示错误原因，直接返回
			// TODO: 错误处理
			return "", errors.UnpaidOrderExists("Unpaid order exists")
		}
	}

	// 通过uuid服务生成uuid
	res, err := osv.UUIDClient.GetUUID(1)
	if err != nil {
		return "", errors.InternalServerError(fmt.Sprintf("Get uuid failed: %v", err))
	}
	orderID = res.ID
	osv.Order.SetID(orderID)
	// 创建订单，与资源库端口（南向网关进行交互）
	err = osv.Port.CreateOrder(osv.Order, siteCode)
	if err != nil {
		return "", errors.InternalServerError(fmt.Sprintf("Create order failed: %v", err))
	}
	return
}

// GetOrderDetail 获取订单详情
// func (osv *OrderService) GetOrderDetail(siteCode string) (order model.Order, err error) {
// 	order, err = osv.Port.GetOrderDetail(osv.Order.GetID(), siteCode)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// // GetOrderList 获取订单列表
// func (osv *OrderService) GetOrderList(args ohs_pl.ListOrderParams) ([]model.Order, int, error) {
// 	// 根据指定id列表获取订单列表
// 	orders, total, err := osv.Port.GetOrderList(args)
// 	if err != nil {
// 		return nil, total, err
// 	}

// 	return orders, total, nil
// }

// UpdateOrderStatus 更新订单状态
func (osv *OrderService) UpdateOrderStatus(siteCode string, status common.StatusType) (err error) {
	// 校验订单是否存在
	_, err = osv.Port.CheckOrderExists(osv.Order.GetID(), siteCode)
	if err != nil {
		return
	}
	// 更新订单状态
	err = osv.Port.UpdateOrderStatus(osv.Order.GetID(), siteCode, status)
	if err != nil {
		return
	}
	return nil
}
