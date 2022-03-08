package clients

// PackageRepository 套餐端口
type PackageRepository interface {
	ValidatePackage(id string) (bool, error)
}
