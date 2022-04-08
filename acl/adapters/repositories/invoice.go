package repositories

import (
	"order-context/acl/adapters/pl"
	"order-context/acl/ports/repositories"
	"order-context/domain/aggregate"
	ohs_pl "order-context/ohs/local/pl"
	"order-context/utils/common"
	"order-context/utils/common/db"
	"sync"
	"time"

	"gorm.io/gorm"
)

// InvoiceAdapter 发票适配器，实现发票端口定义的方法
type InvoiceAdapter struct {
	db *gorm.DB
}

var (
	iOnce sync.Once
	i     repositories.InvoiceRepository
)

// 检查是否实现了接口
var _ repositories.InvoiceRepository = (*InvoiceAdapter)(nil)

// NewInvoiceAdapter 适配器构造方法
func NewInvoiceAdapter() repositories.InvoiceRepository {
	iOnce.Do(func() {
		i = &InvoiceAdapter{
			// 创建数据库引擎
			db: db.NewDBEngine(),
		}
	})
	return i
}

// CreateInvoice 创建发票
func (a *InvoiceAdapter) CreateInvoice(i *aggregate.InvoiceAggregate, siteCode string) error {
	invoice := pl.Invoice{
		ID:        i.Invoice.ID,
		Status:    i.Invoice.Status,
		Code:      common.GenerateNumber(),
		Name:      i.Invoice.Detail.Name,
		Price:     i.Order.Price,
		Path:      i.Invoice.Path,
		OrderID:   i.RootID,
		SiteCode:  siteCode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if res := a.db.Create(&invoice); res.Error != nil {
		return res.Error
	}

	return nil
}

// GetInvoiceDetail 获取发票详情
func (a *InvoiceAdapter) GetInvoiceDetail(invoiceID, siteCode string) (pl.Invoice, error) {
	var invoice pl.Invoice
	if res := a.db.Where("id = ? AND site_code = ?", invoiceID, siteCode).First(&invoice); res.Error != nil {
		return pl.Invoice{}, res.Error
	}

	return invoice, nil
}

// GetInvoiceList 获取发票列表
func (a *InvoiceAdapter) GetInvoiceList(params ohs_pl.ListInvoiceParams) ([]pl.Invoice, int, error) {
	filter := a.db.Table("invoices").Where("order_id = ?", params.OrderID)
	if params.Status != 0 {
		filter = filter.Where("status = ?", params.Status)
	}
	var total int64
	if err := filter.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var invoices []pl.Invoice
	if err := filter.Limit(params.Limit).Offset(params.Offset).Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, int(total), nil
}

// CheckInvoiceExists 检查发票是否存在
func (a *InvoiceAdapter) CheckInvoiceExists(invoiceID, siteCode string) error {
	var invoice pl.Invoice
	if err := a.db.Select("id").Where("id = ? AND site_code = ?", invoiceID, siteCode).First(&invoice).Error; err != nil {
		return err
	}
	return nil
}

// UpdateInvoice 更新发票
func (a *InvoiceAdapter) UpdateInvoice(invoiceID, siteCode string, params ohs_pl.UpdateInvoiceParams) error {
	updateParam := make(map[string]interface{})
	if params.Status != 0 {
		updateParam["status"] = params.Status
	}
	if params.Path != "" {
		updateParam["path"] = params.Path
	}
	updateParam["updated_at"] = time.Now()
	if err := a.db.Model(&pl.Invoice{}).Where("id = ? AND site_code = ?", invoiceID, siteCode).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}
