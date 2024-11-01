package repository

import (
	"github.com/marvelalexius/jones/model"
	"gorm.io/gorm"
)

type (
	NotificationRepository struct {
		db *gorm.DB
	}

	INotificationRepository interface {
		Create(notif model.Notification) error
	}
)

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notif model.Notification) error {
	return r.db.Table("notifications").Create(&notif).Error
}
