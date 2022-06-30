package pl

import "time"

// OrderResponse ...
type OrderResponse struct {
	Status         OrderStatus
	ID             string
	Number         string
	SpaceID        string
	PayID          string
	Price          float32
	PackageVersion string
	PackagePrice   float32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CreateOrderParams ...
type CreateOrderParams struct {
	Price          float32     `json:"price,omitempty,required"`
	Status         OrderStatus `json:"status,omitempty,required"`
	PackagePrice   float32     `json:"package_price,omitempty"`
	PackageVersion string      `json:"package_version,omitempty"`
	SpaceID        string      `json:"space_id,omitempty,mrequired"`
}

// CreateOrderResult ...
type CreateOrderResult struct {
	ID string
}

// UpdateOrderStatusParam ...
type UpdateOrderStatusParam struct {
	ID string `json:"id,omitempty"`
}
