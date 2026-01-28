package models

import (
	"novaardiansyah/simple-pos/pkg/auth"
	"time"
)

type PersonalAccessToken struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	TokenableType string     `json:"tokenable_type"`
	TokenableID   uint       `json:"tokenable_id"`
	Name          string     `json:"name"`
	Token         string     `json:"token"`
	Abilities     string     `json:"abilities"`
	ParentID      *uint      `json:"parent_id"`
	LastUsedAt    *time.Time `json:"last_used_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}

func NewAccessToken(userID uint, duration time.Duration, parentID *uint) (*PersonalAccessToken, string) {
	rawToken, hashedToken := auth.GenerateSecureString(32)
	expiresAt := time.Now().Add(duration)

	name := "refresh_token"
	if parentID != nil {
		name = "auth_token"
	}

	return &PersonalAccessToken{
		Name:          name,
		TokenableID:   userID,
		TokenableType: "App\\Models\\User",
		Token:         hashedToken,
		ParentID:      parentID,
		ExpiresAt:     &expiresAt,
		Abilities:     "[\"*\"]",
	}, rawToken
}
