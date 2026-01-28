/*
 * Project Name: service
 * File: auth_sevice.go
 * Created Date: Tuesday January 27th 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/simple-pos-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"novaardiansyah/simple-pos/internal/dto"
	"novaardiansyah/simple-pos/internal/models"
	"novaardiansyah/simple-pos/internal/repositories"
	"novaardiansyah/simple-pos/pkg/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
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

func (s *authService) Login(c *fiber.Ctx) error {
	data := make(map[string]interface{})

	rules := govalidator.MapData{
		"email":    []string{"required", "email"},
		"password": []string{"required", "min:6"},
	}

	errs := utils.ValidateJSON(c, &data, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	user, err := s.UserRepo.FindByEmail(data["email"].(string))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"].(string)))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	refreshToken, refreshTokenPlain, err := s.generateRefreshToken(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	_, fullToken, err := s.generateAuthToken(user, refreshToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenPlain,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		Path:     "/api/auth/refresh",
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return utils.SuccessResponse(c, "Login successful", dto.LoginResponse{
		Token: fullToken,
	})
}

func (s *authService) ChangePassword(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	user, err := s.UserRepo.FindByID(userId)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User not found")
	}

	data := make(map[string]interface{})

	rules := govalidator.MapData{
		"current_password":          []string{"required", "min:6"},
		"new_password":              []string{"required", "min:6"},
		"new_password_confirmation": []string{"required", "min:6"},
	}

	errs := utils.ValidateJSON(c, &data, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	if data["new_password"] != data["new_password_confirmation"] {
		return utils.ValidationError(c, map[string][]string{
			"new_password": {"Password confirmation does not match"},
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["current_password"].(string)))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(data["new_password"].(string)), 12)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate password")
	}

	hashedPassword := strings.Replace(string(hashed), "$2a$", "$2y$", 1)
	s.UserRepo.UpdatePassword(user.ID, hashedPassword)

	refreshToken, refreshTokenPlain, err := s.generateRefreshToken(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	s.TokenRepo.DeleteByUserID(user.ID)
	_, fullToken, err := s.generateAuthToken(user, refreshToken)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenPlain,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		Path:     "/api/auth/refresh",
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return utils.SuccessResponse(c, "Password changed successfully", dto.LoginResponse{
		Token: fullToken,
	})
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

func (s *authService) generateRefreshToken(user *models.User) (*models.PersonalAccessToken, string, error) {
	length := 40
	bytes := make([]byte, length)
	rand.Read(bytes)

	plainToken := hex.EncodeToString(bytes)[:length]
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	expiration := time.Now().AddDate(0, 0, 7)

	token := models.PersonalAccessToken{
		TokenableType: "App\\Models\\User",
		TokenableID:   user.ID,
		Name:          "refresh_token",
		Token:         hashedToken,
		Abilities:     "[\"*\"]",
		ExpiresAt:     &expiration,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.TokenRepo.Create(&token); err != nil {
		return nil, "", errors.New("token_creation_failed")
	}

	fullToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

	return &token, fullToken, nil
}

func (s *authService) generateAuthToken(user *models.User, refreshToken *models.PersonalAccessToken) (*models.PersonalAccessToken, string, error) {
	length := 40
	bytes := make([]byte, length)
	rand.Read(bytes)

	plainToken := hex.EncodeToString(bytes)[:length]
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	expiration := time.Now().Add(time.Hour)

	token := models.PersonalAccessToken{
		TokenableType: "App\\Models\\User",
		TokenableID:   user.ID,
		Name:          "auth_token",
		Token:         hashedToken,
		Abilities:     "[\"*\"]",
		ExpiresAt:     &expiration,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ParentID:      refreshToken.ID,
	}

	if err := s.TokenRepo.Create(&token); err != nil {
		return nil, "", errors.New("token_creation_failed")
	}

	fullToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

	return &token, fullToken, nil
}
