package model

import "time"

// NotificationRepository :nodoc:
type NotificationRepository interface {
	Create(notif Notification) error
}

// Notification :nodoc:
type Notification struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Content   string     `json:"content"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `gorm:"<-:create" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// MatchMessage :nodoc:
type MatchMessage struct {
	Type    string `json:"type"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func NewMatchNotification(target, content string) Notification {
	return Notification{
		UserID:  target,
		Content: content,
	}
}
