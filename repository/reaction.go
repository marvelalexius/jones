package repository

import (
	"context"

	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/marvelalexius/jones/utils/str"
	"gorm.io/gorm"
)

type (
	ReactionRepository struct {
		db *gorm.DB
	}

	IReactionRepository interface {
		FindMatch(ctx context.Context, userID, matchedUserID string) (model.Reaction, error)
		FindSwiped(ctx context.Context, userID string) (reactions []model.Reaction, err error)
		FindSwipeCount(ctx context.Context, userID string) (int64, error)
		Create(ctx context.Context, reaction model.Reaction) error
		Update(ctx context.Context, reaction *model.Reaction) error
	}
)

func NewReactionRepository(db *gorm.DB) IReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) FindSwiped(ctx context.Context, userID string) (reactions []model.Reaction, err error) {
	err = r.db.Table("reactions").Where("user_id = ?", userID).Find(&reactions).Error
	if err != nil {
		logger.Errorln(ctx, "failed to find swiped", err)

		return reactions, err
	}

	return reactions, nil
}

func (r *ReactionRepository) FindMatch(ctx context.Context, userID, matchedUserID string) (reactions model.Reaction, err error) {
	err = r.db.Table("reactions").Where("user_id = ?", userID).Where("matched_user_id = ?", matchedUserID).Where("type = ?", model.ReactionLike).First(&reactions).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorln(ctx, "failed to find match", err)

		return reactions, err
	}

	return reactions, nil
}

func (r *ReactionRepository) FindSwipeCount(ctx context.Context, userID string) (int64, error) {
	var count int64

	nowEarliest, nowLatest := str.GetTodayTimeRange()
	err := r.db.Table("reactions").Where("user_id = ?", userID).Where("created_at between ? and ?", nowEarliest, nowLatest).Count(&count).Error
	if err != nil {
		logger.Errorln(ctx, "failed to find swipe count", err)

		return count, err
	}

	return count, nil
}

func (r *ReactionRepository) Create(ctx context.Context, reaction model.Reaction) error {
	return r.db.Table("reactions").Create(&reaction).Error
}

func (r *ReactionRepository) Update(ctx context.Context, reaction *model.Reaction) error {
	return r.db.Table("reactions").Where("id = ?", reaction.ID).Updates(&reaction).Error
}
