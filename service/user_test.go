package service

import (
	"context"
	"errors"
	"testing"

	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/mocks"
	"github.com/marvelalexius/jones/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name          string
		input         model.LoginUser
		mockSetup     func(*mocks.IUserRepository)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "successful login",
			input: model.LoginUser{
				Email:    "test@example.com",
				Password: "testtest",
			},
			mockSetup: func(ur *mocks.IUserRepository) {
				user := &model.User{
					ID:       "user123",
					Email:    "test@example.com",
					Password: "$2a$10$JLS6wpwU9KXbApMoIZz9tee54U.x7efqOTwkySALnFwmmK7nOyeLi", // Pre-hashed password
				}
				ur.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectedUser: &model.User{
				ID:       "user123",
				Email:    "test@example.com",
				Password: "$2a$10$JLS6wpwU9KXbApMoIZz9tee54U.x7efqOTwkySALnFwmmK7nOyeLi",
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			input: model.LoginUser{
				Email:    "nonexistent@example.com",
				Password: "testtest",
			},
			mockSetup: func(ur *mocks.IUserRepository) {
				ur.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "invalid password",
			input: model.LoginUser{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(ur *mocks.IUserRepository) {
				user := &model.User{
					ID:       "user123",
					Email:    "test@example.com",
					Password: "$2a$10$JLS6wpwU9KXbApMoIZz9tee54U.x7efqOTwkySALnFwmmK7nOyeLi",
				}
				ur.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectedUser:  nil,
			expectedError: errors.New("invalid username or password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.IUserRepository)
			reactionRepo := new(mocks.IReactionRepository)
			tt.mockSetup(userRepo)

			service := NewUserService(&config.Config{
				App: config.App{
					Secret:             "some-secret-key",
					RefreshTokenSecret: "some-refresh-token-secret",
				},
			}, userRepo, reactionRepo)
			user, err := service.Login(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name          string
		input         *model.RegisterUser
		mockSetup     func(*mocks.IUserRepository)
		expectedError error
	}{
		{
			name: "successful registration",
			input: &model.RegisterUser{
				Email:       "new@example.com",
				Password:    "validpassword",
				Name:        "New User",
				DateOfBirth: "2000-01-01",
			},
			mockSetup: func(ur *mocks.IUserRepository) {
				ur.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, gorm.ErrRecordNotFound)
				ur.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "user already exists",
			input: &model.RegisterUser{
				Email:    "existing@example.com",
				Password: "password",
				Name:     "Existing User",
			},
			mockSetup: func(ur *mocks.IUserRepository) {
				existingUser := &model.User{
					ID:    "existing123",
					Email: "existing@example.com",
				}
				ur.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			expectedError: errors.New("user already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.IUserRepository)
			reactionRepo := new(mocks.IReactionRepository)
			tt.mockSetup(userRepo)

			service := NewUserService(&config.Config{}, userRepo, reactionRepo)
			user, err := service.Register(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, user.ID)
				assert.Equal(t, tt.input.Email, user.Email)
				assert.NotEqual(t, tt.input.Password, user.Password) // Password should be hashed
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_FindAll(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		mockSetup     func(*mocks.IUserRepository, *mocks.IReactionRepository)
		expectedUsers []model.User
		expectedTotal int64
		expectedError error
	}{
		{
			name:   "successful find all",
			userID: "user123",
			mockSetup: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository) {
				loggedInUser := &model.User{
					ID:         "user123",
					Preference: "female",
				}
				ur.On("FindByID", mock.Anything, "user123").Return(loggedInUser, nil)

				swiped := []model.Reaction{
					{MatchedUserID: "user456"},
					{MatchedUserID: "user789"},
				}
				rr.On("FindSwiped", mock.Anything, "user123").Return(swiped, nil)

				users := []model.User{
					{ID: "user111", Name: "User 1"},
					{ID: "user222", Name: "User 2"},
				}
				ur.On("FindAll", mock.Anything, mock.Anything, "female").Return(users, int64(2), nil)
			},
			expectedUsers: []model.User{
				{ID: "user111", Name: "User 1"},
				{ID: "user222", Name: "User 2"},
			},
			expectedTotal: 2,
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			mockSetup: func(ur *mocks.IUserRepository, rr *mocks.IReactionRepository) {
				ur.On("FindByID", mock.Anything, "nonexistent").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedUsers: []model.User{},
			expectedTotal: 0,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.IUserRepository)
			reactionRepo := new(mocks.IReactionRepository)
			tt.mockSetup(userRepo, reactionRepo)

			service := NewUserService(&config.Config{}, userRepo, reactionRepo)
			users, total, err := service.FindAll(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUsers, users)
				assert.Equal(t, tt.expectedTotal, total)
			}
			userRepo.AssertExpectations(t)
			reactionRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_RefreshAuthToken(t *testing.T) {
	tests := []struct {
		name            string
		refreshToken    string
		mockSetup       func(*mocks.IUserRepository)
		expectedToken   string
		expectedRefresh string
		expectedError   error
	}{
		{
			name:         "successful token refresh",
			refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIwMUpCS1c2WEYzS0pLS0g0NFlOODAwVktIVyIsImV4cCI6MTczMTA4Nzg0NH0.vMDTG1XImQFRIYZJQTesvMJKKHpY8psVZ_1nVjg8QHg",
			mockSetup: func(ur *mocks.IUserRepository) {
				user := &model.User{
					ID:    "user123",
					Email: "test@example.com",
				}
				ur.On("FindByID", mock.Anything, mock.Anything).Return(user, nil)
			},
			expectedError: nil,
		},
		{
			name:          "invalid refresh token",
			refreshToken:  "invalid.token",
			mockSetup:     func(ur *mocks.IUserRepository) {},
			expectedError: errors.New("invalid token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(mocks.IUserRepository)
			reactionRepo := new(mocks.IReactionRepository)
			tt.mockSetup(userRepo)

			config := &config.Config{
				App: config.App{
					Secret:             "testsecret",
					RefreshTokenSecret: "testrefreshsecret",
				},
			}

			service := NewUserService(config, userRepo, reactionRepo)
			token, refresh, err := service.RefreshAuthToken(context.Background(), tt.refreshToken)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, refresh)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GenerateAuthTokens(t *testing.T) {
	config := &config.Config{
		App: config.App{
			Secret:             "testsecret",
			RefreshTokenSecret: "testrefreshsecret",
		},
	}

	user := &model.User{
		ID:    "user123",
		Email: "test@example.com",
	}

	t.Run("successful token generation", func(t *testing.T) {
		userRepo := new(mocks.IUserRepository)
		reactionRepo := new(mocks.IReactionRepository)
		service := NewUserService(config, userRepo, reactionRepo)

		token, refresh, err := service.GenerateAuthTokens(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, refresh)
	})
}
