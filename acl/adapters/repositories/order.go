package repositories

import (
	"fmt"
	"order-service/acl/adapters/pl"
	"order-service/acl/ports/repositories"
	"order-service/common"
	"order-service/domain/aggregate"
	"sync"
	"time"

	"gorm.io/gorm"
)

// OrderAdapter 订单适配器，实现订单端口定义的方法
type OrderAdapter struct {
	db *gorm.DB
}

var (
	oOnce sync.Once
	oa    repositories.OrderRepository
)

// 检查是否实现了接口
var _ repositories.OrderRepository = (*OrderAdapter)(nil)

// NewOrderAdapter 订单适配器构造方法
func NewOrderAdapter() repositories.OrderRepository {
	oOnce.Do(func() {
		oa = &OrderAdapter{
			// 创建数据库引擎
			db: common.NewDBEngine(),
		}
	})
	return oa
}

// CreateOrder 创建订单
func (a *OrderAdapter) CreateOrder(root *aggregate.AggregateRoot, siteCode string) error {
	order := pl.Order{
		ID:             root.Order.ID,
		Status:         root.Order.Status,
		Number:         common.GenerateNumber(),
		SpaceID:        root.Space.ID,
		PayID:          root.Payment.Voucher,
		Price:          root.Order.Price,
		PackageVersion: root.Package.Version,
		PackagePrice:   root.Package.Price,
		SiteCode:       siteCode,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := a.db.Create(&order).Error; err != nil {
		fmt.Print(err)
		return err
	}
	return nil

}

// GetOrderDetail 获取订单详情
func (a *OrderAdapter) GetOrderDetail(orderID, siteCode string) (pl.Order, error) {
	var order pl.Order
	if res := a.db.Where("id = ? AND site_code = ?", orderID, siteCode).First(&order); res.Error != nil {
		return pl.Order{}, res.Error
	}

	return order, nil
}

// GetOrderList 获取订单列表
func (a *OrderAdapter) GetOrderList(params common.ListOrderParams) ([]pl.Order, error) {
	filter := a.db.Table("orders").Where("space_id = ?", params.SpaceID)
	if params.Status != 0 {
		filter = filter.Where("status = ?", params.Status)
	}
	var total int64
	if err := filter.Count(&total).Error; err != nil {
		return nil, err
	}
	var orders []pl.Order
	if err := filter.Limit(params.Limit).Offset(params.Offset).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// CheckOrderExists 检查订单是否存在
func (a *OrderAdapter) CheckOrderExists(orderID, siteCode string) (bool, error) {
	var order pl.Order
	if err := a.db.Select("id").Where("id = ? AND site_code = ?", orderID, siteCode).First(&order).Error; err != nil {
		return false, err
	}
	return true, nil
}

// UpdateOrderStatus 更新订单
func (a *OrderAdapter) UpdateOrderStatus(orderID, siteCode string, status common.StatusType) error {
	if err := a.db.Model(&pl.Order{}).Where("id = ? AND site_code = ?", orderID, siteCode).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error; err != nil {
		return err
	}

	return nil
}
