package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/model"
	stripePkg "github.com/marvelalexius/jones/pkg/stripe"
	"github.com/marvelalexius/jones/repository"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/stripe/stripe-go/v81"
	"gorm.io/gorm"
)

type (
	SubscriptionService struct {
		Conf             *config.Config
		StripeClient     stripePkg.IStripeClient
		UserRepo         repository.IUserRepository
		SubscriptionRepo repository.ISubscriptionRepository
	}

	ISubscriptionService interface {
		Subscribe(ctx context.Context, userID string, req model.SubscriptionRequest) (*stripe.CheckoutSession, error)
		CustomerPortal(ctx context.Context, userID string) (*stripe.BillingPortalSession, error)
		HandleInvoicePaid(ctx context.Context, invoice *stripe.Invoice) error
		HandleInvoicePaymentFailed(ctx context.Context, customerEmail string) error
		HandleSubscriptionDeleted(ctx context.Context, stripeSubscription *stripe.Subscription) error
	}
)

func NewSubscriptionService(conf *config.Config, stripeClient stripePkg.IStripeClient, userRepo repository.IUserRepository, subscriptionRepo repository.ISubscriptionRepository) ISubscriptionService {
	return &SubscriptionService{Conf: conf, StripeClient: stripeClient, UserRepo: userRepo, SubscriptionRepo: subscriptionRepo}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, userID string, req model.SubscriptionRequest) (*stripe.CheckoutSession, error) {
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "error finding user by ID", err)

	}

	if user.StripeCustomerID == "" {
		customer, err := s.StripeClient.CreateCustomer(ctx, user.Email, user.Name)
		if err != nil {
			logger.Errorln(ctx, "error requesting create customer to stripe", err)

			return nil, err
		}

		user.StripeCustomerID = customer.ID
		user, err = s.UserRepo.Update(user)
		if err != nil {
			logger.Errorln(ctx, "error updating user", err)

			return nil, err
		}
	}

	checkoutSession, err := s.StripeClient.CreateCheckoutSession(ctx, user.StripeCustomerID, req.PlanID)
	if err != nil {
		logger.Errorln(ctx, "error requesting create checkout session to stripe", err)

		return nil, err
	}

	return checkoutSession, nil
}

func (s *SubscriptionService) CustomerPortal(ctx context.Context, userID string) (*stripe.BillingPortalSession, error) {
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "error finding user by ID", err)

		return nil, err
	}

	if user.StripeCustomerID == "" {
		return nil, errors.New("customer not found")
	}

	billingPortalSession, err := s.StripeClient.CreateBillingPortalSession(ctx, user.StripeCustomerID)
	if err != nil {
		logger.Errorln(ctx, "error requesting create billing portal session to stripe", err)

		return nil, err
	}

	return billingPortalSession, nil
}

func (s *SubscriptionService) HandleInvoicePaid(ctx context.Context, invoice *stripe.Invoice) error {
	stripeProduct := invoice.Lines.Data[0]

	user, err := s.UserRepo.FindByEmail(ctx, invoice.CustomerEmail)
	if err != nil {
		logger.Errorln(ctx, "error finding user by email", err)

		return err
	}

	subscription, err := s.SubscriptionRepo.FindByUserID(ctx, user.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorln(ctx, "error finding subscription by user ID", err)

		return err
	}

	if subscription != nil {
		if subscription.StripeSubscriptionID == invoice.Subscription.ID {
			subscription.ExpiredAt = time.Unix(int64(stripeProduct.Period.End), 0)

			err = s.SubscriptionRepo.Update(ctx, subscription)
			if err != nil {
				logger.Errorln(ctx, "error updating subscription", err)

				return err
			}

			return nil
		} else {
			now := time.Now()
			subscription.CanceledAt = &now

			err = s.SubscriptionRepo.Update(ctx, subscription)
			if err != nil {
				logger.Errorln(ctx, "error updating subscription", err)

				return err
			}
		}
	}

	plan, err := s.SubscriptionRepo.FindPlanByProductID(ctx, stripeProduct.Plan.Product.ID)
	if err != nil {
		logger.Errorln(ctx, "error finding plan by product ID", err)

		return err
	}

	periodStart := time.Unix(int64(stripeProduct.Period.Start), 0)
	periodEnd := time.Unix(int64(stripeProduct.Period.End), 0)
	newSubscription := model.Subscription{
		UserID:               user.ID,
		StripeSubscriptionID: invoice.Subscription.ID,
		PlanID:               plan.ID,
		StartedAt:            periodStart,
		ExpiredAt:            periodEnd,
	}

	err = s.SubscriptionRepo.Create(ctx, newSubscription)
	if err != nil {
		logger.Errorln(ctx, "error creating subscription", err)

		return err
	}

	return nil
}

func (s *SubscriptionService) HandleInvoicePaymentFailed(ctx context.Context, customerEmail string) error {
	user, err := s.UserRepo.FindByEmail(ctx, customerEmail)
	if err != nil {
		logger.Errorln(ctx, "error finding user by email", err)

		return err
	}

	subscription, err := s.SubscriptionRepo.FindByUserID(ctx, user.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorln(ctx, "error finding subscription by user ID", err)

		return err
	}

	if subscription != nil {
		// TODO: notify user that payment failed and their subscription is still active until expiration date
	}

	// TODO: notify user that payment failed

	return nil
}

func (s *SubscriptionService) HandleSubscriptionDeleted(ctx context.Context, stripeSubscription *stripe.Subscription) error {
	subscription, err := s.SubscriptionRepo.FindByStripeSubscriptionID(ctx, stripeSubscription.ID)
	if err != nil {
		logger.Errorln(ctx, "error finding subscription by stripe subscription ID", err)

		return err
	}

	user, err := s.UserRepo.FindByID(ctx, subscription.UserID)
	if err != nil {
		logger.Errorln(ctx, "error finding user by ID", err)

		return err
	}

	canceledAt := time.Unix(int64(stripeSubscription.CanceledAt), 0)
	subscription.CanceledAt = &canceledAt

	err = s.SubscriptionRepo.Update(ctx, subscription)
	if err != nil {
		logger.Errorln(ctx, "error updating subscription", err)

		return err
	}

	// TODO: notify user that their subscription has been canceled

	fmt.Println(user)
	return nil
}
