package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	SubscriptionPlanFree  = "free"
	SubscriptionPlanBasic = "basic"
	SubscriptionPlanPro   = "pro"

	SubscriptionStatusActive   = "active"
	SubscriptionStatusExpired  = "expired"
	SubscriptionStatusCanceled = "canceled"
)

type SubscriptionRequest struct {
	PlanID string `json:"plan_id"`
}

type SubscriptionPlan struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Price         decimal.Decimal `json:"price"`
	Features      []string        `json:"features"`
	StripePriceID string          `json:"-"`
	CreatedAt     time.Time       `gorm:"<-:create" json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
}

type Subscription struct {
	ID                   string           `json:"id"`
	UserID               string           `json:"user_id"`
	PlanID               string           `json:"plan_id"`
	StripeSubscriptionID string           `json:"-"`
	StartedAt            time.Time        `gorm:"<-:create" json:"started_at"`
	ExpiredAt            time.Time        `gorm:"<-:create" json:"expired_at"`
	CanceledAt           *time.Time       `json:"canceled_at"`
	CreatedAt            time.Time        `gorm:"<-:create" json:"created_at"`
	UpdatedAt            *time.Time       `json:"updated_at"`
	SubscriptionPlan     SubscriptionPlan `gorm:"->" json:"subscription_plan"`
}
