package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	ReactionLike    = "LIKE"
	ReactionDislike = "PASS"
)

type Reaction struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	MatchedUserID string     `json:"matched_user_id"`
	Type          string     `json:"type"`
	MatchedAt     *time.Time `json:"matched_at"`
	CreatedAt     time.Time  `gorm:"<-:create" json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type ReactionRequest struct {
	UserID        string `json:"user_id"`
	MatchedUserID string `json:"matched_user_id"`
	Type          string `json:"type"`
}

func (r *ReactionRequest) ToReactionModel() Reaction {
	return Reaction{
		ID:            ulid.Make().String(),
		UserID:        r.UserID,
		MatchedUserID: r.MatchedUserID,
		Type:          r.Type,
	}
}
