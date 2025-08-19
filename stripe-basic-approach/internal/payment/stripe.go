package payment

type StripeProcessor struct{}

func NewStripeProcessor() PaymentProcessor {
	return &StripeProcessor{}
}

func (s *StripeProcessor) CreatePayment() error { return nil }
