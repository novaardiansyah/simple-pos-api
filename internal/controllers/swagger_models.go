package controllers

import "time"

type UserSwagger struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	HasAllowNotification *bool     `json:"has_allow_notification"`
	NotificationToken    *string   `json:"notification_token,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DeletedAt            *string   `json:"deleted_at,omitempty"`
}
