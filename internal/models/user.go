package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	Code                 string         `gorm:"size:255" json:"code"`
	Name                 string         `gorm:"size:255;not null" json:"name"`
	Email                string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password             string         `gorm:"size:255;not null" json:"-"`
	HasAllowNotification *bool          `gorm:"default:false" json:"has_allow_notification"`
	NotificationToken    *string        `gorm:"size:255" json:"notification_token,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
}

// TableName specifies table name
func (User) TableName() string {
	return "users"
}
