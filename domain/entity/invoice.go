package entity

import (
	"order-context/domain/vo"
	"order-context/utils/common"
)

// Invoice 实体
type Invoice struct {
	// ID 发票ID
	ID string
	// 发票状态
	Status common.InvoiceStatusType
	// 发票保存路径
	Path string
	// 发票详情
	Detail vo.InvoiceDetail
}
