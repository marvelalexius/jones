package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/repository"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	ReactionService struct {
		UserRepo         repository.IUserRepository
		ReactionRepo     repository.IReactionRepository
		SubscriptionRepo repository.ISubscriptionRepository
		NotificationRepo repository.INotificationRepository
	}

	IReactionService interface {
		Swipe(ctx context.Context, reaction model.ReactionRequest) (model.Reaction, error)
		SeeLikes(ctx context.Context, userID string) ([]model.Reaction, error)
	}
)

func NewReactionService(userRepo repository.IUserRepository, reactionRepo repository.IReactionRepository, subscriptionRepo repository.ISubscriptionRepository, notificationRepo repository.INotificationRepository) IReactionService {
	return &ReactionService{UserRepo: userRepo, ReactionRepo: reactionRepo, SubscriptionRepo: subscriptionRepo, NotificationRepo: notificationRepo}
}

func (s *ReactionService) Swipe(ctx context.Context, req model.ReactionRequest) (model.Reaction, error) {
	subscribed, err := s.SubscriptionRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		logger.Errorln(ctx, "failed to check subscription", err)

		return model.Reaction{}, errors.New("failed to check subscription")
	}

	if subscribed.ID == "" {
		count, err := s.ReactionRepo.FindSwipeCount(ctx, req.UserID)
		if err != nil {
			logger.Errorln(ctx, "failed to check swipe count", err)

			return model.Reaction{}, errors.New("failed to check swipe count")
		}

		if count >= 10 {
			logger.Errorln(ctx, "cannot swipe more than 10 times", err)

			return model.Reaction{}, errors.New("cannot swipe more than 10 times. please try again tomorrow")
		}
	}

	hasSwiped, err := s.ReactionRepo.HasSwiped(ctx, req.UserID, req.MatchedUserID)
	if err != nil {
		logger.Errorln(ctx, "failed to check if user has swiped", err)

		return model.Reaction{}, errors.New("failed to check if user has swiped")
	}

	if hasSwiped.ID != "" {
		logger.Errorln(ctx, "user has already swiped")

		return model.Reaction{}, errors.New("user has already swiped")
	}

	reaction := req.ToReactionModel()
	matched, err := s.ReactionRepo.FindMatch(ctx, req.MatchedUserID, req.UserID)
	if err != nil {
		logger.Errorln(ctx, "failed to find match", err)

		return model.Reaction{}, errors.New("failed to find match")
	}

	if matched.ID == "" {
		err = s.ReactionRepo.Create(ctx, reaction)
		if err != nil {
			logger.Errorln(ctx, "failed to create reaction", err)

			return model.Reaction{}, errors.New("failed to create reaction")
		}

		return reaction, nil
	}

	now := time.Now()
	reaction.MatchedAt = &now

	matched.MatchedAt = &now
	matched.UpdatedAt = &now

	err = s.ReactionRepo.Update(ctx, &matched)
	if err != nil {
		logger.Errorln(ctx, "failed to update reaction", err)

		return model.Reaction{}, errors.New("failed to update reaction")
	}

	err = s.ReactionRepo.Create(ctx, reaction)
	if err != nil {
		logger.Errorln(ctx, "failed to create reaction", err)

		return model.Reaction{}, errors.New("failed to create reaction")
	}

	// send notification to swipe
	s.sendMatchNotification(reaction)

	// send notification to matched
	s.sendMatchNotification(matched)

	return reaction, nil
}

func (s *ReactionService) SeeLikes(ctx context.Context, userID string) ([]model.Reaction, error) {
	subscribed, err := s.SubscriptionRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "failed to check subscription", err)

		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("you are not subscribed to any plan")
		}

		return nil, errors.New("failed to check subscription")
	}

	plan, err := s.SubscriptionRepo.FindPlanByID(ctx, subscribed.PlanID)
	if err != nil {
		logger.Errorln(ctx, "failed to find plan by ID", err)

		return nil, errors.New("failed to find plan")
	}

	if plan.Name != model.SubscriptionPlanPro {
		return nil, errors.New("you are not a pro user")
	}

	reactions, err := s.ReactionRepo.FindLikes(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "failed to find likes", err)

		return nil, errors.New("failed to find likes")
	}

	return reactions, nil
}

func (u *ReactionService) sendMatchNotification(reaction model.Reaction) {
	content := model.MatchMessage{
		Type:    model.ReactionLike,
		UserID:  reaction.MatchedUserID,
		Message: "Congratulations! You matched",
	}

	b, err := json.Marshal(content)
	if err != nil {
		logrus.WithField("reaction", reaction).Error(err.Error())
	}

	notification := model.Notification{
		ID:        ulid.Make().String(),
		UserID:    reaction.UserID,
		Content:   string(b),
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	err = u.NotificationRepo.Create(notification)
	if err != nil {
		logrus.WithField("reaction", reaction).Error(err.Error())
	}
}
