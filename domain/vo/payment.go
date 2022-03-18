package vo

// Payment 支付值对象
type Payment struct {
	// 支付单号
	Voucher string
	// 支付来源
	Source string
	// 币种
	Currency string
	// 金额
	Total float32
}
