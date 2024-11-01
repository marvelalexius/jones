package repository

import (
	"context"
	"time"

	"github.com/marvelalexius/jones/model"
	"gorm.io/gorm"
)

type (
	SubscriptionRepository struct {
		db *gorm.DB
	}

	ISubscriptionRepository interface {
		Create(ctx context.Context, subscription model.Subscription) error
		Update(ctx context.Context, subscription *model.Subscription) error
		FindByID(ctx context.Context, id string) (*model.Subscription, error)
		FindByStripeSubscriptionID(ctx context.Context, id string) (*model.Subscription, error)
		FindAllPlan(ctx context.Context) ([]model.SubscriptionPlan, error)
		FindPlanByProductID(ctx context.Context, id string) (*model.SubscriptionPlan, error)
		FindAll(ctx context.Context) ([]model.Subscription, error)
		FindByUserID(ctx context.Context, userID string) (*model.Subscription, error)
	}
)

func NewSubscriptionRepository(db *gorm.DB) ISubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription model.Subscription) error {
	return r.db.Table("subscriptions").Create(&subscription).Error
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	return r.db.Table("subscriptions").Where("id = ?", subscription.ID).Updates(&subscription).Error
}

func (r *SubscriptionRepository) FindByID(ctx context.Context, id string) (*model.Subscription, error) {
	var subscription model.Subscription

	if err := r.db.First(&subscription, id).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) FindAllPlan(ctx context.Context) ([]model.SubscriptionPlan, error) {
	var subscriptionPlans []model.SubscriptionPlan

	if err := r.db.Find(&subscriptionPlans).Error; err != nil {
		return nil, err
	}

	return subscriptionPlans, nil
}

func (r *SubscriptionRepository) FindPlanByProductID(ctx context.Context, id string) (*model.SubscriptionPlan, error) {
	var subscriptionPlan model.SubscriptionPlan

	if err := r.db.Where("stripe_product_id = ?", id).First(&subscriptionPlan).Error; err != nil {
		return nil, err
	}

	return &subscriptionPlan, nil
}

func (r *SubscriptionRepository) FindAll(ctx context.Context) ([]model.Subscription, error) {
	var subscriptions []model.Subscription

	if err := r.db.Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) FindByUserID(ctx context.Context, userID string) (*model.Subscription, error) {
	var subscription model.Subscription

	if err := r.db.Where("user_id = ?", userID).Where("canceled_at is null").Or("expired_at <= ?", time.Now().Format("2006-01-02")).Find(&subscription).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) FindByStripeSubscriptionID(ctx context.Context, id string) (*model.Subscription, error) {
	var subscription model.Subscription

	if err := r.db.Where("stripe_subscription_id = ?", id).First(&subscription).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}
