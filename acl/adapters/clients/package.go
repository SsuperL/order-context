package clients

// PackageAdapter 接收消息契约模型
type PackageAdapter struct {
}

// ValidatePackage 校验套餐有效性
func (p *PackageAdapter) ValidatePackage(id string) (bool, error) {
	return false, nil
}
