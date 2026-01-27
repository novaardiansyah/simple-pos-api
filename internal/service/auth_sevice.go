package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"novaardiansyah/simple-pos/internal/models"
	"novaardiansyah/simple-pos/internal/repositories"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (string, error)
	ChangePassword(user *models.User, currentPassword, newPassword string) (string, error)
	UpdateProfile(user *models.User, name, email string) error
}

type authService struct {
	UserRepo  *repositories.UserRepository
	TokenRepo *repositories.PersonalAccessTokenRepository
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		UserRepo:  repositories.NewUserRepository(db),
		TokenRepo: repositories.NewPersonalAccessTokenRepository(db),
	}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid_credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid_credentials")
	}

	return s.generateAuthToken(user, 7)
}

func (s *authService) ChangePassword(user *models.User, currentPassword, newPassword string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return "", errors.New("current_password_incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return "", err
	}

	hashedPassword := strings.Replace(string(hashed), "$2a$", "$2y$", 1)
	s.UserRepo.UpdatePassword(user.ID, hashedPassword)

	s.TokenRepo.DeleteByUserID(user.ID)
	newToken, err := s.generateAuthToken(user, 7)

	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (s *authService) UpdateProfile(user *models.User, name, email string) error {
	if email == "" {
		email = user.Email
	}

	if user.Email != email {
		if _, err := s.UserRepo.FindByEmail(email); err == nil {
			return errors.New("email_already_used")
		}
	}

	updateFields := map[string]interface{}{
		"name":  name,
		"email": email,
	}

	if err := s.UserRepo.UpdateFields(user.ID, updateFields); err != nil {
		return err
	}

	return nil
}

func (s *authService) generateAuthToken(user *models.User, expireDays int) (string, error) {
	length := 40
	bytes := make([]byte, length)
	rand.Read(bytes)

	plainToken := hex.EncodeToString(bytes)[:length]
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	expiration := time.Now().AddDate(0, 0, expireDays)

	token := models.PersonalAccessToken{
		TokenableType: "App\\Models\\User",
		TokenableID:   user.ID,
		Name:          "auth_token",
		Token:         hashedToken,
		Abilities:     "[\"*\"]",
		ExpiresAt:     &expiration,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.TokenRepo.Create(&token); err != nil {
		return "", errors.New("token_creation_failed")
	}

	fullToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

	return fullToken, nil
}
