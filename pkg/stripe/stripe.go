package stripe

import (
	"context"

	"github.com/marvelalexius/jones/utils/logger"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

type (
	StripeClient struct {
		Client        *client.API
		Secret        string
		WebhookSecret string
	}

	IStripeClient interface {
		CreateCustomer(ctx context.Context, email string, name string) (*stripe.Customer, error)
		// CreatePaymentMethod(ctx context.Context, customerID string, cardNumber string, cardCVC string, cardExpMonth string, cardExpYear string) (string, error)
		CreateCheckoutSession(ctx context.Context, customerID string, planID string) (*stripe.CheckoutSession, error)
		CreateBillingPortalSession(ctx context.Context, customerID string) (*stripe.BillingPortalSession, error)
	}
)

func NewStripeClient(secret string, webhookSecret string) IStripeClient {
	sc := &client.API{}
	sc.Init(secret, nil)

	return &StripeClient{
		Client: sc,
	}
}

func (c *StripeClient) CreateCustomer(ctx context.Context, email string, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}

	result, err := c.Client.Customers.New(params)
	if err != nil {
		logger.Errorln(ctx, "failed to create customer", err)

		return nil, err
	}

	return result, nil
}

func (c *StripeClient) CreateCheckoutSession(ctx context.Context, customerID string, planID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer:   stripe.String(customerID),
		SuccessURL: stripe.String("https://example.com/success"),
		CancelURL:  stripe.String("https://example.com/cancel"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(planID),
				Quantity: stripe.Int64(1),
			},
		},
	}

	session, err := c.Client.CheckoutSessions.New(params)
	if err != nil {
		logger.Errorln(ctx, "failed to create checkout session", err)

		return nil, err
	}

	return session, err
}

func (c *StripeClient) CreateBillingPortalSession(ctx context.Context, customerID string) (*stripe.BillingPortalSession, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer: stripe.String(customerID),
	}

	session, err := c.Client.BillingPortalSessions.New(params)
	if err != nil {
		logger.Errorln(ctx, "failed to create billing portal session", err)

		return nil, err
	}

	return session, err
}
