package service

import (
	"context"
	"errors"
	"time"

	"github.com/marvelalexius/jones/config"
	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/repository"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/marvelalexius/jones/utils/str"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type (
	UserService struct {
		Config       *config.Config
		UserRepo     repository.IUserRepository
		ReactionRepo repository.IReactionRepository
	}

	IUserService interface {
		Login(ctx context.Context, req model.LoginUser) (*model.User, error)
		Register(ctx context.Context, user *model.RegisterUser) (*model.User, error)
		FindAll(ctx context.Context, userID string) (users []model.User, total int64, err error)
		RefreshAuthToken(ctx context.Context, refreshToken string) (string, string, error)
		GenerateAuthTokens(user *model.User) (string, string, error)
	}
)

func NewUserService(config *config.Config, userRepo repository.IUserRepository, reactionRepo repository.IReactionRepository) IUserService {
	return &UserService{Config: config, UserRepo: userRepo, ReactionRepo: reactionRepo}
}

func (s *UserService) Login(ctx context.Context, req model.LoginUser) (*model.User, error) {
	user, err := s.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		logger.Errorln(ctx, "failed to find user", err)

		return nil, err
	}

	err = user.CheckPassword(req.Password)
	if err != nil {
		logger.Errorln(ctx, "failed to check password", err)

		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterUser) (*model.User, error) {
	userExists, err := s.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorln(ctx, "failed to find user", err)

		return &model.User{}, err
	}

	if userExists != nil {
		return &model.User{}, errors.New("user already exists")
	}

	user := req.ToUserModel()
	user.ID = ulid.Make().String()

	if err := user.HashPassword(req.Password); err != nil {
		logger.Errorln(ctx, "failed to hash password", err)

		return &model.User{}, err
	}

	if len(req.Images) >= 1 {
		user.NewImageFromRequest(req.Images)
	}

	if err := s.UserRepo.Create(user); err != nil {
		logger.Errorln(ctx, "failed to create user", err)

		return &model.User{}, err
	}

	return user, nil
}

func (s *UserService) RefreshAuthToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := str.ParseJWT(refreshToken, s.Config.App.Secret)
	if err != nil {
		logger.Errorln(ctx, "failed to parse refresh token", err)

		return "", "", err
	}

	user, err := s.UserRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		logger.Errorln(ctx, "failed to find user", err)

		return "", "", err
	}

	token, refreshToken, err := s.GenerateAuthTokens(user)
	if err != nil {
		logger.Errorln(ctx, "failed to generate auth tokens", err)

		return "", "", err
	}

	return token, refreshToken, nil
}

func (s *UserService) FindAll(ctx context.Context, userID string) (users []model.User, total int64, err error) {
	loggedInUser, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "failed to get logged in user", err)

		return []model.User{}, 0, err
	}

	swiped, err := s.ReactionRepo.FindSwiped(ctx, userID)
	if err != nil {
		logger.Errorln(ctx, "failed to find swiped", err)

		return []model.User{}, 0, err
	}

	userIDs := []string{loggedInUser.ID}
	for _, swipedUser := range swiped {
		userIDs = append(userIDs, swipedUser.MatchedUserID)
	}

	users, total, err = s.UserRepo.FindAll(ctx, userIDs)
	if err != nil {
		logger.Errorln(ctx, "failed to find users", err)

		return []model.User{}, 0, err
	}

	return users, total, err
}

func (a *UserService) GenerateAuthTokens(user *model.User) (string, string, error) {
	token, err := str.GenerateJWT(user.ID, time.Now().Add(24*time.Hour), a.Config.App.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := str.GenerateJWT(user.ID, time.Now().Add(7*24*time.Hour), a.Config.App.RefreshTokenSecret)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
