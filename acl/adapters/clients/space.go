package clients

// SpaceAdapter 空间适配器
type SpaceAdapter struct {
}

// ValidateSpace 校验空间有效性
func (u *SpaceAdapter) ValidateSpace(SpaceID string) (bool, error) {
	return false, nil
}
