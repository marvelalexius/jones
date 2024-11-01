package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/mocks"
	"github.com/marvelalexius/jones/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v76"
	"gorm.io/gorm"
)

func TestSubscribe(t *testing.T) {
	ctx := context.Background()
	mockStripeClient := new(mocks.IStripeClient)
	mockUserRepo := new(mocks.IUserRepository)
	mockSubscriptionRepo := new(mocks.ISubscriptionRepository)

	conf := &config.Config{}
	subscriptionService := NewSubscriptionService(conf, mockStripeClient, mockUserRepo, mockSubscriptionRepo)

	tests := []struct {
		name          string
		userID        string
		request       model.SubscriptionRequest
		setupMocks    func()
		expectedURL   string
		expectedError error
	}{
		{
			name:   "Success - New Customer",
			userID: "user123",
			request: model.SubscriptionRequest{
				PlanID: 1,
			},
			setupMocks: func() {
				mockUserRepo.On("FindByID", ctx, "user123").Return(&model.User{
					ID:    "user123",
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)

				mockStripeClient.On("CreateCustomer", ctx, "test@example.com", "Test User").Return(&stripe.Customer{
					ID: "cus_123",
				}, nil)

				mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(&model.User{
					ID:               "user123",
					StripeCustomerID: "cus_123",
				}, nil)

				mockSubscriptionRepo.On("FindPlanByID", ctx, 1).Return(&model.SubscriptionPlan{
					ID:            1,
					StripePriceID: "prod_123",
				}, nil)

				mockStripeClient.On("CreateCheckoutSession", ctx, "cus_123", "prod_123").Return(&stripe.CheckoutSession{
					URL: "https://checkout.stripe.com/session",
				}, nil)
			},
			expectedURL:   "https://checkout.stripe.com/session",
			expectedError: nil,
		},
		{
			name:   "Success - Existing Customer",
			userID: "user456",
			request: model.SubscriptionRequest{
				PlanID: 1,
			},
			setupMocks: func() {
				mockUserRepo.ExpectedCalls = nil
				mockStripeClient.ExpectedCalls = nil
				mockSubscriptionRepo.ExpectedCalls = nil

				mockUserRepo.On("FindByID", ctx, "user456").Return(&model.User{
					ID:               "user456",
					StripeCustomerID: "cus_456",
					Email:            "test@example.com",
					Name:             "Test User",
				}, nil)

				mockSubscriptionRepo.On("FindPlanByID", ctx, 1).Return(&model.SubscriptionPlan{
					ID:            1,
					StripePriceID: "prod_456",
				}, nil)

				mockStripeClient.On("CreateCheckoutSession", ctx, "cus_456", "prod_456").Return(&stripe.CheckoutSession{
					URL: "https://checkout.stripe.com/session2",
				}, nil)
			},
			expectedURL:   "https://checkout.stripe.com/session2",
			expectedError: nil,
		},
		{
			name:   "Error - User Not Found",
			userID: "nonexistent",
			request: model.SubscriptionRequest{
				PlanID: 1,
			},
			setupMocks: func() {
				mockUserRepo.On("FindByID", ctx, "nonexistent").Return(&model.User{}, gorm.ErrRecordNotFound)
			},
			expectedURL:   "",
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			url, err := subscriptionService.Subscribe(ctx, tt.userID, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}
		})
	}
}

func TestCustomerPortal(t *testing.T) {
	ctx := context.Background()
	mockStripeClient := new(mocks.IStripeClient)
	mockUserRepo := new(mocks.IUserRepository)
	mockSubscriptionRepo := new(mocks.ISubscriptionRepository)

	conf := &config.Config{}
	subscriptionService := NewSubscriptionService(conf, mockStripeClient, mockUserRepo, mockSubscriptionRepo)

	tests := []struct {
		name          string
		userID        string
		setupMocks    func()
		expectedURL   string
		expectedError error
	}{
		{
			name:   "Success",
			userID: "user123",
			setupMocks: func() {
				mockUserRepo.On("FindByID", ctx, "user123").Return(&model.User{
					ID:               "user123",
					StripeCustomerID: "cus_123",
				}, nil)

				mockStripeClient.On("CreateBillingPortalSession", ctx, "cus_123").Return(&stripe.BillingPortalSession{
					URL: "https://billing.stripe.com/session",
				}, nil)
			},
			expectedURL:   "https://billing.stripe.com/session",
			expectedError: nil,
		},
		{
			name:   "Error - No Customer ID",
			userID: "user456",
			setupMocks: func() {
				mockUserRepo.On("FindByID", ctx, "user456").Return(&model.User{
					ID: "user456",
				}, nil)
			},
			expectedURL:   "",
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			url, err := subscriptionService.CustomerPortal(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}
		})
	}
}

func TestHandleInvoicePaid(t *testing.T) {
	ctx := context.Background()
	mockStripeClient := new(mocks.IStripeClient)
	mockUserRepo := new(mocks.IUserRepository)
	mockSubscriptionRepo := new(mocks.ISubscriptionRepository)

	conf := &config.Config{}
	subscriptionService := NewSubscriptionService(conf, mockStripeClient, mockUserRepo, mockSubscriptionRepo)

	testTime := time.Now()
	periodEnd := int64(time.Now().Add(30 * 24 * time.Hour).Unix())

	tests := []struct {
		name          string
		invoice       *stripe.Invoice
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success - New Subscription",
			invoice: &stripe.Invoice{
				CustomerEmail: "test@example.com",
				Subscription:  &stripe.Subscription{ID: "sub_123"},
				Lines: &stripe.InvoiceLineItemList{
					Data: []*stripe.InvoiceLineItem{
						{
							Price: &stripe.Price{ID: "price_123"},
							Period: &stripe.Period{
								Start: int64(testTime.Unix()),
								End:   periodEnd,
							},
						},
					},
				},
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", ctx, "test@example.com").Return(&model.User{
					ID: "user123",
				}, nil)

				mockSubscriptionRepo.On("FindByUserID", ctx, "user123").Return(nil, gorm.ErrRecordNotFound)

				mockSubscriptionRepo.On("FindPlanByProductID", ctx, "price_123").Return(&model.SubscriptionPlan{
					ID: 1,
				}, nil)

				mockSubscriptionRepo.On("Create", ctx, mock.AnythingOfType("model.Subscription")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Success - Existing Subscription Update",
			invoice: &stripe.Invoice{
				CustomerEmail: "test@example.com",
				Subscription:  &stripe.Subscription{ID: "sub_123"},
				Lines: &stripe.InvoiceLineItemList{
					Data: []*stripe.InvoiceLineItem{
						{
							Price: &stripe.Price{ID: "price_123"},
							Period: &stripe.Period{
								Start: int64(testTime.Unix()),
								End:   periodEnd,
							},
						},
					},
				},
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", ctx, "test@example.com").Return(&model.User{
					ID: "user123",
				}, nil)

				mockSubscriptionRepo.On("FindByUserID", ctx, "user123").Return(&model.Subscription{
					ID:                   "sub_existing",
					StripeSubscriptionID: "sub_123",
					PlanID:               1,
				}, nil)

				mockSubscriptionRepo.On("FindPlanByProductID", ctx, "price_123").Return(&model.SubscriptionPlan{
					ID: 1,
				}, nil)

				mockSubscriptionRepo.On("Update", ctx, mock.AnythingOfType("*model.Subscription")).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := subscriptionService.HandleInvoicePaid(ctx, tt.invoice)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleSubscriptionUpdated(t *testing.T) {
	ctx := context.Background()
	mockStripeClient := new(mocks.IStripeClient)
	mockUserRepo := new(mocks.IUserRepository)
	mockSubscriptionRepo := new(mocks.ISubscriptionRepository)

	conf := &config.Config{}
	subscriptionService := NewSubscriptionService(conf, mockStripeClient, mockUserRepo, mockSubscriptionRepo)

	periodEnd := int64(time.Now().Add(30 * 24 * time.Hour).Unix())

	tests := []struct {
		name               string
		stripeSubscription *stripe.Subscription
		setupMocks         func()
		expectedError      error
	}{
		{
			name: "Success - Active Subscription Update",
			stripeSubscription: &stripe.Subscription{
				ID:               "sub_123",
				CurrentPeriodEnd: periodEnd,
				Items: &stripe.SubscriptionItemList{
					Data: []*stripe.SubscriptionItem{
						{
							Plan: &stripe.Plan{ID: "plan_123"},
						},
					},
				},
			},
			setupMocks: func() {
				mockSubscriptionRepo.On("FindPlanByProductID", ctx, "plan_123").Return(&model.SubscriptionPlan{
					ID: 1,
				}, nil)

				mockSubscriptionRepo.On("FindByStripeSubscriptionIDAndPlanID", ctx, "sub_123", 1).Return(&model.Subscription{
					ID:     "sub_existing",
					UserID: "user123",
				}, nil)

				mockUserRepo.On("FindByID", ctx, "user123").Return(&model.User{
					ID: "user123",
				}, nil)

				mockSubscriptionRepo.On("Update", ctx, mock.AnythingOfType("*model.Subscription")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Success - Canceled Subscription Update",
			stripeSubscription: &stripe.Subscription{
				ID:               "sub_123",
				CurrentPeriodEnd: periodEnd,
				CanceledAt:       time.Now().Unix(),
				Items: &stripe.SubscriptionItemList{
					Data: []*stripe.SubscriptionItem{
						{
							Plan: &stripe.Plan{ID: "plan_123"},
						},
					},
				},
			},
			setupMocks: func() {
				mockSubscriptionRepo.On("FindPlanByProductID", ctx, "plan_123").Return(&model.SubscriptionPlan{
					ID: 1,
				}, nil)

				mockSubscriptionRepo.On("FindByStripeSubscriptionIDAndPlanID", ctx, "sub_123", 1).Return(&model.Subscription{
					ID:     "sub_existing",
					UserID: "user123",
				}, nil)

				mockUserRepo.On("FindByID", ctx, "user123").Return(&model.User{
					ID: "user123",
				}, nil)

				mockSubscriptionRepo.On("Update", ctx, mock.AnythingOfType("*model.Subscription")).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := subscriptionService.HandleSubscriptionUpdated(ctx, tt.stripeSubscription)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleInvoicePaymentFailed(t *testing.T) {
	ctx := context.Background()
	mockStripeClient := new(mocks.IStripeClient)
	mockUserRepo := new(mocks.IUserRepository)
	mockSubscriptionRepo := new(mocks.ISubscriptionRepository)

	conf := &config.Config{}
	subscriptionService := NewSubscriptionService(conf, mockStripeClient, mockUserRepo, mockSubscriptionRepo)

	tests := []struct {
		name          string
		customerEmail string
		setupMocks    func()
		expectedError error
	}{
		{
			name:          "Success - With Active Subscription",
			customerEmail: "test@example.com",
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", ctx, "test@example.com").Return(&model.User{
					ID: "user123",
				}, nil)

				mockSubscriptionRepo.On("FindByUserID", ctx, "user123").Return(&model.Subscription{
					ID:     "sub_existing",
					UserID: "user123",
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:          "Success - No Active Subscription",
			customerEmail: "test@example.com",
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", ctx, "test@example.com").Return(&model.User{
					ID: "user123",
				}, nil)

				mockSubscriptionRepo.On("FindByUserID", ctx, "user123").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: nil,
		},
		{
			name:          "Error - User Not Found",
			customerEmail: "nonexistent@example.com",
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", ctx, "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := subscriptionService.HandleInvoicePaymentFailed(ctx, tt.customerEmail)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
