package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/model"
	stripePkg "github.com/marvelalexius/jones/pkg/stripe"
	"github.com/marvelalexius/jones/repository"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/oklog/ulid/v2"
	"github.com/stripe/stripe-go/v76"
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
		Subscribe(ctx context.Context, userID string, req model.SubscriptionRequest) (string, error)
		CustomerPortal(ctx context.Context, userID string) (string, error)
		HandleInvoicePaid(ctx context.Context, invoice *stripe.Invoice) error
		HandleInvoicePaymentFailed(ctx context.Context, customerEmail string) error
		HandleSubscriptionUpdated(ctx context.Context, stripeSubscription *stripe.Subscription) error
		HandleSubscriptionDeleted(ctx context.Context, stripeSubscription *stripe.Subscription) error
	}
)

func NewSubscriptionService(conf *config.Config, stripeClient stripePkg.IStripeClient, userRepo repository.IUserRepository, subscriptionRepo repository.ISubscriptionRepository) ISubscriptionService {
	return &SubscriptionService{Conf: conf, StripeClient: stripeClient, UserRepo: userRepo, SubscriptionRepo: subscriptionRepo}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, userID string, req model.SubscriptionRequest) (string, error) {
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "error finding user by ID", err)

		return "", err
	}

	plan, err := s.SubscriptionRepo.FindPlanByID(ctx, req.PlanID)
	if err != nil {
		logger.Errorln(ctx, "error finding plan by ID", err)

		return "", err
	}

	if s.Conf.FeatureFlag.EnableStripe && plan.StripePriceID != "" {
		if user.StripeCustomerID == "" {
			customer, err := s.StripeClient.CreateCustomer(ctx, user.Email, user.Name)
			if err != nil {
				logger.Errorln(ctx, "error requesting create customer to stripe", err)

				return "", err
			}

			user.StripeCustomerID = customer.ID
			user, err = s.UserRepo.Update(user)
			if err != nil {
				logger.Errorln(ctx, "error updating user", err)

				return "", err
			}
		}

		checkoutSession, err := s.StripeClient.CreateCheckoutSession(ctx, user.StripeCustomerID, plan.StripePriceID)
		if err != nil {
			logger.Errorln(ctx, "error requesting create checkout session to stripe", err)

			return "", err
		}

		return checkoutSession.URL, nil
	} else {
		newSubscription := model.Subscription{
			ID:        ulid.Make().String(),
			UserID:    user.ID,
			PlanID:    plan.ID,
			StartedAt: time.Now(),
			ExpiredAt: time.Now().Add(30 * 24 * time.Hour),
		}

		err = s.SubscriptionRepo.Create(ctx, newSubscription)
		if err != nil {
			logger.Errorln(ctx, "error creating subscription", err)

			return "", err
		}

		return newSubscription.ID, nil
	}
}

func (s *SubscriptionService) CustomerPortal(ctx context.Context, userID string) (string, error) {
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "error finding user by ID", err)

		return "", err
	}

	if user.StripeCustomerID == "" {
		return "", errors.New("customer not found")
	}

	billingPortalSession, err := s.StripeClient.CreateBillingPortalSession(ctx, user.StripeCustomerID)
	if err != nil {
		logger.Errorln(ctx, "error requesting create billing portal session to stripe", err)

		return "", err
	}

	return billingPortalSession.URL, nil
}

func (s *SubscriptionService) HandleInvoicePaid(ctx context.Context, invoice *stripe.Invoice) error {
	/*
	*	By default, Stripe always sends invoice.paid events when the invoice is paid, whether it's for a new subscription, switching plans, or renewal.
	*
	*	Checking the invoice lines data below is for a workaround when user switch plans using billing portal.
	*	Stripe returns the same subscription ID but added a new price data in the invoice lines data for the new plan.
	*
	 */
	stripeProduct := invoice.Lines.Data[0]
	if len(invoice.Lines.Data) > 1 {
		stripeProduct = invoice.Lines.Data[1]
	}

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

	paidPlan, err := s.SubscriptionRepo.FindPlanByProductID(ctx, stripeProduct.Price.ID)
	if err != nil {
		logger.Errorln(ctx, "error finding plan by product ID", err)

		return err
	}

	/*
	*
	*	This logic checks if the user has already subscribed before. If they have, we need to update their subscription.
	*
	*	If the latest subscription's Stripe subscription ID is equal to the Stripe subscription ID in the invoice, and the subscribed plan ID is exactly the same as the new plan ID, we can update the expired at to use the latest one and return immediately.
	*
	*	If the latest subscription's Stripe subscription ID is not equal to the Stripe subscription ID in the invoice, we need to update the subscription to canceled, and create a new one based on the new subscribed plan.
	*
	 */
	if subscription != nil {
		if subscription.StripeSubscriptionID == invoice.Subscription.ID && subscription.PlanID == paidPlan.ID {
			subscription.ExpiredAt = time.Unix(int64(stripeProduct.Period.End), 0)

			err = s.SubscriptionRepo.Update(ctx, subscription)
			if err != nil {
				logger.Errorln(ctx, "error updating subscription", err)

				return err
			}

			return nil
		} else {
			now := time.Now()
			subscription.CanceledAt = sql.NullTime{Time: now, Valid: true}

			err = s.SubscriptionRepo.Update(ctx, subscription)
			if err != nil {
				logger.Errorln(ctx, "error updating subscription", err)

				return err
			}
		}
	}

	periodStart := time.Unix(int64(stripeProduct.Period.Start), 0)
	periodEnd := time.Unix(int64(stripeProduct.Period.End), 0)
	newSubscription := model.Subscription{
		ID:                   ulid.Make().String(),
		UserID:               user.ID,
		StripeSubscriptionID: invoice.Subscription.ID,
		PlanID:               paidPlan.ID,
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
		logger.Infoln(ctx, "your payment has failed. your subscription is still active until expiration date")
		return nil
	}

	// TODO: notify user that payment failed
	logger.Infoln(ctx, "your payment has failed. Not to worry, your credit card has not been charged. please try again, or contact support")

	return nil
}

func (s *SubscriptionService) HandleSubscriptionUpdated(ctx context.Context, stripeSubscription *stripe.Subscription) error {
	stripeProduct := stripeSubscription.Items.Data[0]

	updatedPlan, err := s.SubscriptionRepo.FindPlanByProductID(ctx, stripeProduct.Plan.ID)
	if err != nil {
		logger.Errorln(ctx, "error finding plan by product ID", err)

		return err
	}

	subscription, err := s.SubscriptionRepo.FindByStripeSubscriptionIDAndPlanID(ctx, stripeSubscription.ID, updatedPlan.ID)
	if err != nil {
		logger.Errorln(ctx, "error finding subscription by stripe subscription ID", err)

		return err
	}

	user, err := s.UserRepo.FindByID(ctx, subscription.UserID)
	if err != nil {
		logger.Errorln(ctx, "error finding user by ID", err)

		return err
	}

	/*
	*	Stripe's billing portal allows users to cancel their subscription at any time and renew at any time.
	*	We need to handle the cancellation of the subscription here.
	*
	*	https://stripe.com/docs/billing/billing-portal
	*
	*	This logic checks if the subscription has been canceled and updates the subscription accordingly.
	*	If the canceledAtis not 0, it means the subscription has been canceled.
	*	If the canceledAt is 0, it means the subscription has not been canceled or switched plans.
	 */
	if stripeSubscription.CanceledAt != 0 {
		canceledAt := time.Unix(int64(stripeSubscription.CanceledAt), 0)
		subscription.CanceledAt = sql.NullTime{Time: canceledAt, Valid: true}
	} else {
		subscription.CanceledAt = sql.NullTime{Time: time.Now(), Valid: false}
	}

	subscription.PlanID = updatedPlan.ID
	subscription.ExpiredAt = time.Unix(int64(stripeSubscription.CurrentPeriodEnd), 0)

	err = s.SubscriptionRepo.Update(ctx, subscription)
	if err != nil {
		logger.Errorln(ctx, "error updating subscription", err)

		return err
	}

	// TODO: notify user that their subscription has been canceled
	if stripeSubscription.CanceledAt != 0 {
		logger.Infoln(ctx, fmt.Sprintf("your subscription has been canceled. You still have access to your account until %s", subscription.ExpiredAt.Format("2006-01-02")))
	}

	fmt.Println(user)
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
	subscription.CanceledAt = sql.NullTime{Time: canceledAt, Valid: true}

	err = s.SubscriptionRepo.Update(ctx, subscription)
	if err != nil {
		logger.Errorln(ctx, "error updating subscription", err)

		return err
	}

	// TODO: notify user that their subscription has been canceled
	logger.Infoln(ctx, fmt.Sprintf("your subscription has been canceled. You still have access to your account until %s", subscription.ExpiredAt.Format("2006-01-02")))

	fmt.Println(user)
	return nil
}
