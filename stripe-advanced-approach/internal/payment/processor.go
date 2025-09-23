package payment

import (
	"context"

	"github.com/google/uuid"
)

type PaymentProcessor interface {
	SetupProducts(context.Context, *SetupProductsReq) (*SetupProductsResp, error)
	SetupSubscription(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error)
	GetProducts(ctx context.Context) (*ProductListResponse, error)
	CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error)
	SaveCard(ctx context.Context, customerId string) (string, error)
	CreatePaymentIntent(ctx context.Context, amount int64, customerId string) (*CreatePaymentIntentResponse, error)
	PurchaseProduct(ctx context.Context, req *PurchaseProductRequest) (*StripePurchaseResponse, error)
	SubscribeToProduct(ctx context.Context, req *SubscribeRequest) (*SubscribeResponse, error)
}
