package clients

type PaymentRepository interface {
	Pay() error
	CallBack() error
}
