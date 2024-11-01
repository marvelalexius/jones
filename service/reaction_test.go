package service

import (
	"context"
	"errors"
	"testing"

	"github.com/marvelalexius/jones/mocks"
	"github.com/marvelalexius/jones/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestReactionService_Swipe(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		request       model.ReactionRequest
		setupMocks    func(*mocks.IUserRepository, *mocks.IReactionRepository, *mocks.ISubscriptionRepository, *mocks.INotificationRepository)
		expectedError error
	}{
		{
			name: "Error - User Subscription Not Found",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, gorm.ErrRecordNotFound)
			},
			expectedError: errors.New("failed to check subscription"),
		},
		{
			name: "Error - Find Swipe Count",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), gorm.ErrRecordNotFound).Once()
			},
			expectedError: errors.New("failed to check swipe count"),
		},
		{
			name: "Error - Failed to Check If User Has Swiped",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil).Once()
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, gorm.ErrRecordNotFound).Once()
			},
			expectedError: errors.New("failed to check if user has swiped"),
		},
		{
			name: "Error - Find Match",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil).Once()
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, nil).Once()
				rr.On("FindMatch", mock.Anything, "user2", "user1").Return(model.Reaction{}, gorm.ErrRecordNotFound).Once()
			},
			expectedError: errors.New("failed to find match"),
		},
		{
			name: "Error - Failed to Create Reaction",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil).Once()
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, nil).Once()
				rr.On("FindMatch", mock.Anything, "user2", "user1").Return(model.Reaction{}, nil).Once()
				rr.On("Create", mock.Anything, mock.AnythingOfType("model.Reaction")).Return(gorm.ErrRecordNotFound).Once()
			},
			expectedError: errors.New("failed to create reaction"),
		},
		{
			name: "Error - Update Matched User's Data",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil).Once()
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, nil).Once()

				matchedReaction := model.Reaction{
					ID:            "reaction1",
					UserID:        "user2",
					MatchedUserID: "user1",
					Type:          model.ReactionLike,
				}
				rr.On("FindMatch", mock.Anything, "user2", "user1").Return(matchedReaction, nil)
				rr.On("Update", mock.Anything, mock.AnythingOfType("*model.Reaction")).Return(gorm.ErrRecordNotFound)
			},
			expectedError: errors.New("failed to update reaction"),
		},
		{
			name: "Error - Update Matched User's Data",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil).Once()
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, nil).Once()

				matchedReaction := model.Reaction{
					ID:            "reaction1",
					UserID:        "user2",
					MatchedUserID: "user1",
					Type:          model.ReactionLike,
				}
				rr.On("FindMatch", mock.Anything, "user2", "user1").Return(matchedReaction, nil)
				rr.On("Update", mock.Anything, mock.AnythingOfType("*model.Reaction")).Return(nil)
				rr.On("Create", mock.Anything, mock.AnythingOfType("model.Reaction")).Return(errors.New("failed to create reaction")).Once()
			},
			expectedError: errors.New("failed to create reaction"),
		},
		{
			name: "Success - First Swipe No Match",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil).Once()
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, nil).Once()
				rr.On("FindMatch", mock.Anything, "user2", "user1").Return(model.Reaction{}, nil).Once()
				rr.On("Create", mock.Anything, mock.AnythingOfType("model.Reaction")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Success - Match Found",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil).Once()
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, nil)

				matchedReaction := model.Reaction{
					ID:            "reaction1",
					UserID:        "user2",
					MatchedUserID: "user1",
					Type:          model.ReactionLike,
				}
				rr.On("FindMatch", mock.Anything, "user2", "user1").Return(matchedReaction, nil)
				rr.On("Update", mock.Anything, mock.AnythingOfType("*model.Reaction")).Return(nil)
				rr.On("Create", mock.Anything, mock.AnythingOfType("model.Reaction")).Return(nil)
				nr.On("Create", mock.AnythingOfType("model.Notification")).Return(nil)
				nr.On("Create", mock.AnythingOfType("model.Notification")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Success - Pro User Can Swipe more than 10 times",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.ExpectedCalls = nil
				rr.ExpectedCalls = nil
				ur.ExpectedCalls = nil
				nr.ExpectedCalls = nil

				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{
					ID: "sub1",
				}, nil)
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{}, nil).Once()
				rr.On("FindMatch", mock.Anything, "user2", "user1").Return(model.Reaction{}, nil).Once()
				rr.On("Create", mock.Anything, mock.AnythingOfType("model.Reaction")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Error - Already Swiped",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(0), nil)
				rr.On("HasSwiped", mock.Anything, "user1", "user2").Return(model.Reaction{ID: "existing"}, nil)
			},
			expectedError: errors.New("user has already swiped"),
		},
		{
			name: "Error - Swipe Limit Exceeded",
			request: model.ReactionRequest{
				UserID:        "user1",
				MatchedUserID: "user2",
				Type:          model.ReactionLike,
			},
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, nil)
				rr.On("FindSwipeCount", mock.Anything, "user1").Return(int64(10), nil)
			},
			expectedError: errors.New("cannot swipe more than 10 times. please try again tomorrow"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			userRepo := new(mocks.IUserRepository)
			reactionRepo := new(mocks.IReactionRepository)
			subscriptionRepo := new(mocks.ISubscriptionRepository)
			notificationRepo := new(mocks.INotificationRepository)

			// Setup mocks
			tt.setupMocks(userRepo, reactionRepo, subscriptionRepo, notificationRepo)

			// Create service
			service := NewReactionService(userRepo, reactionRepo, subscriptionRepo, notificationRepo)

			// Execute
			reaction, err := service.Swipe(ctx, tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, reaction.ID)
				assert.Equal(t, tt.request.UserID, reaction.UserID)
				assert.Equal(t, tt.request.MatchedUserID, reaction.MatchedUserID)
				assert.Equal(t, tt.request.Type, reaction.Type)
			}

			// Verify all mocks
			userRepo.AssertExpectations(t)
			reactionRepo.AssertExpectations(t)
			subscriptionRepo.AssertExpectations(t)
			notificationRepo.AssertExpectations(t)
		})
	}
}

func TestReactionService_SeeLikes(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*mocks.IUserRepository, *mocks.IReactionRepository, *mocks.ISubscriptionRepository, *mocks.INotificationRepository)
		expectedError error
	}{
		{
			name:   "Success - Pro User",
			userID: "user1",
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{
					ID:     "sub1",
					PlanID: 2,
				}, nil)
				sr.On("FindPlanByID", mock.Anything, 2).Return(&model.SubscriptionPlan{
					ID:   2,
					Name: model.SubscriptionPlanPro,
				}, nil)
				rr.On("FindLikes", mock.Anything, "user1").Return([]model.Reaction{
					{
						ID:            "reaction1",
						UserID:        "user2",
						MatchedUserID: "user1",
						Type:          model.ReactionLike,
					},
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:   "Error - General Error",
			userID: "user1",
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, gorm.ErrInvalidDB)
			},
			expectedError: errors.New("failed to check subscription"),
		},
		{
			name:   "Error - Not Subscribed",
			userID: "user1",
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{}, gorm.ErrRecordNotFound)
			},
			expectedError: errors.New("you are not subscribed to any plan"),
		},
		{
			name:   "Error - Not Pro User",
			userID: "user1",
			setupMocks: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository, sr *mocks.ISubscriptionRepository, nr *mocks.INotificationRepository) {
				sr.On("FindByUserID", mock.Anything, "user1").Return(&model.Subscription{
					ID:     "sub1",
					PlanID: 1,
				}, nil)
				sr.On("FindPlanByID", mock.Anything, 1).Return(&model.SubscriptionPlan{
					ID:   1,
					Name: model.SubscriptionPlanBasic,
				}, nil)
			},
			expectedError: errors.New("you are not a pro user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			userRepo := new(mocks.IUserRepository)
			reactionRepo := new(mocks.IReactionRepository)
			subscriptionRepo := new(mocks.ISubscriptionRepository)
			notificationRepo := new(mocks.INotificationRepository)

			// Setup mocks
			tt.setupMocks(userRepo, reactionRepo, subscriptionRepo, notificationRepo)

			// Create service
			service := NewReactionService(userRepo, reactionRepo, subscriptionRepo, notificationRepo)

			// Execute
			reactions, err := service.SeeLikes(ctx, tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, reactions)
			}

			// Verify all mocks
			userRepo.AssertExpectations(t)
			reactionRepo.AssertExpectations(t)
			subscriptionRepo.AssertExpectations(t)
			notificationRepo.AssertExpectations(t)
		})
	}
}
