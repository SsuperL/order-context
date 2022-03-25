package pl

import (
	"order-context/common"
	"time"
)

// Order 数据库模型
type Order struct {
	// 订单ID
	ID string `gorm:"column:id;primary_key:not null;comment:订单id" json:"id"`
	// 订单状态
	Status common.StatusType `gorm:"column:status;not null;comment:订单状态"`
	// 订单号
	Number string `gorm:"column:number;not null;comment:订单号"`
	// 空间ID
	SpaceID string `gorm:"column:space_id;index;not null;comment:空间ID"`
	// 支付单号ID
	PayID string `gorm:"column:pay_id;comment:支付单号"`
	// 订单总价
	Price float32 `gorm:"column:price;not null;comment:订单总价"`
	// 套餐版本
	PackageVersion string `gorm:"column:package_version;not null;comment:套餐版本"`
	// 套餐价格
	PackagePrice float32 `gorm:"column:package_price;not null;comment:套餐价格"`
	// 自定义域名site-code
	SiteCode string `gorm:"column:site_code;not null;comment:site-code"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at;not null;comment:创建时间"`
	// 变更时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null;comment:更新时间"`
}

// Invoice 发票数据库模型
type Invoice struct {
	// 发票ID
	ID string `gorm:"column:id;primary_key:not null;comment:发票id" json:"id"`
	// 订单ID
	OrderID string `gorm:"column:order_id;index;not null;comment:订单id" json:"order_id"`
	// 订单总价
	Price float32 `gorm:"column:price;not null;comment:订单总价"`
	// 发票抬头
	Name string `gorm:"column:name;not null;comment:发票抬头" json:"name"`
	// 发票税号
	Code string `gorm:"column:code;not null;comment:发票税号" json:"code"`
	// 发票状态
	Status common.InvoiceStatusType `gorm:"column:status;not null;comment:发票状态" json:"status"`
	// 发票存放路径
	Path string `gorm:"column:path;not null;comment:发票存放路径" json:"path"`
	// 自定义域名site-code
	SiteCode string `gorm:"column:site_code;not null;comment:site-code"`
	// 创建时间
	CreatedAt time.Time `gorm:"column:created_at;not null;comment:创建时间"`
	// 变更时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null;comment:更新时间"`
}

// TableName return tablename of order
func (Order) TableName() string {
	return "orders"
}

// TableName return tablename of invoice
func (Invoice) TableName() string {
	return "invoices"
}
