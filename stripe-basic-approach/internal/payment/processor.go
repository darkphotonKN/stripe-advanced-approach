package payment

type PaymentProcessor interface {
	CreatePayment() error
}
