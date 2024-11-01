package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/repository"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/sirupsen/logrus"
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
	}
)

func NewReactionService(userRepo repository.IUserRepository, reactionRepo repository.IReactionRepository, subscriptionRepo repository.ISubscriptionRepository, notificationRepo repository.INotificationRepository) IReactionService {
	return &ReactionService{UserRepo: userRepo, ReactionRepo: reactionRepo, SubscriptionRepo: subscriptionRepo, NotificationRepo: notificationRepo}
}

func (s *ReactionService) Swipe(ctx context.Context, req model.ReactionRequest) (model.Reaction, error) {
	subscribed, err := s.SubscriptionRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		logger.Errorln(ctx, "failed to check subscription", err)

		return model.Reaction{}, err
	}

	if subscribed == nil {
		count, err := s.ReactionRepo.FindSwipeCount(ctx, req.UserID)
		if err != nil {
			logger.Errorln(ctx, "failed to find swipe count", err)

			return model.Reaction{}, err
		}

		if count >= 10 {
			logger.Errorln(ctx, "cannot swipe more than 10 times", err)

			return model.Reaction{}, err
		}
	}

	reaction := req.ToReactionModel()
	matched, err := s.ReactionRepo.FindMatch(ctx, req.MatchedUserID, req.UserID)
	if err != nil {
		logger.Errorln(ctx, "failed to find match", err)

		return model.Reaction{}, err
	}

	if matched.ID != "" {
		err = s.ReactionRepo.Create(ctx, reaction)
		if err != nil {
			logger.Errorln(ctx, "failed to create reaction", err)

			return model.Reaction{}, err
		}

		return reaction, nil
	}

	now := time.Now()
	reaction.MatchedAt = &now
	matched.MatchedAt = &now

	err = s.ReactionRepo.Update(ctx, &matched)
	if err != nil {
		logger.Errorln(ctx, "failed to update reaction", err)

		return model.Reaction{}, err
	}

	err = s.ReactionRepo.Create(ctx, reaction)
	if err != nil {
		logger.Errorln(ctx, "failed to create reaction", err)

		return model.Reaction{}, err
	}

	// send notification to swipe
	go s.sendMatchNotification(reaction)

	// send notification to matched
	go s.sendMatchNotification(matched)

	return reaction, nil
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
