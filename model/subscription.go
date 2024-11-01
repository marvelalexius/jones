package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

const (
	SubscriptionPlanBasic = "BASIC"
	SubscriptionPlanPro   = "PRO"

	SubscriptionStatusActive   = "active"
	SubscriptionStatusExpired  = "expired"
	SubscriptionStatusCanceled = "canceled"
)

var SubscriptionFeatures = map[string][]string{
	"BASIC": {
		"unlimited_likes",
	},
	"PRO": {
		"unlimited_likes",
		"see_likes",
	},
}

type SubscriptionRequest struct {
	PlanID int `json:"plan_id" binding:"required,numeric"`
}

type SubscriptionPlan struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Price         decimal.Decimal `json:"price"`
	Features      pq.StringArray  `gorm:"type:text[]" json:"features"`
	StripePriceID string          `json:"-"`
	CreatedAt     time.Time       `gorm:"<-:create" json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
}

type Subscription struct {
	ID                   string       `json:"id"`
	UserID               string       `json:"user_id"`
	PlanID               int          `json:"plan_id"`
	StripeSubscriptionID string       `json:"-"`
	StartedAt            time.Time    `json:"started_at"`
	ExpiredAt            time.Time    `json:"expired_at"`
	CanceledAt           sql.NullTime `json:"canceled_at"`
	CreatedAt            time.Time    `gorm:"<-:create" json:"created_at"`
	UpdatedAt            *time.Time   `json:"updated_at"`
}
