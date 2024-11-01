package repository

import (
	"context"

	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/utils/logger"
	"gorm.io/gorm"
)

type (
	UserRepository struct {
		db *gorm.DB
	}

	IUserRepository interface {
		FindAll(ctx context.Context, userIds []string, preference string) (users []model.User, total int64, err error)
		FindByID(ctx context.Context, id string) (*model.User, error)
		FindByEmail(ctx context.Context, email string) (*model.User, error)
		FindByStripeCustomerID(ctx context.Context, id string) (*model.User, error)
		Create(user *model.User) error
		Update(user *model.User) (*model.User, error)
	}
)

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll(ctx context.Context, userIds []string, preference string) (users []model.User, total int64, err error) {
	q := r.db.Table("users").Not("id in (?)", userIds)

	if val, ok := model.SupportedPreference[preference]; ok {
		if val == "BOTH" {
			q = q.Where("(gender = ? OR gender = ?)", model.GenderMale, model.GenderFemale)
		} else {
			q = q.Where("gender = ?", val)
		}
	}

	err = q.Count(&total).Error
	if err != nil {
		logger.Errorln(ctx, "failed to count users", err)

		return users, total, err
	}

	q = q.Order("created_at desc")

	err = q.Preload("Images").Find(&users).Error
	if err != nil {
		logger.Errorln(ctx, "failed to find users", err)

		return users, total, err
	}

	return users, total, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := r.db.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByStripeCustomerID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.Model(&model.User{}).Where("stripe_customer_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	if err := r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
