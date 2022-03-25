package clients

import "order-context/acl/adapters/pl"

// UUIDClient 客户端端口，调用UUID服务
type UUIDClient interface {
	// GetUUID获取uuid
	GetUUID(int) (pl.UUIDRes, error)
}
