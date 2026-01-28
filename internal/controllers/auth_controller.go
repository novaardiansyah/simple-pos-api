/*
 * Project Name: controllers
 * File: auth_controller.go
 * Created Date: Tuesday January 27th 2026
 *
 * Author: Nova Ardiansyah admin@novaardiansyah.id
 * Website: https://novaardiansyah.id
 * MIT License: https://github.com/novaardiansyah/simple-pos-api/blob/main/LICENSE
 *
 * Copyright (c) 2026 Nova Ardiansyah, Org
 */

package controllers

import (
	"novaardiansyah/simple-pos/internal/dto"
	"novaardiansyah/simple-pos/internal/models"
	"novaardiansyah/simple-pos/internal/repositories"
	"novaardiansyah/simple-pos/internal/service"
	"novaardiansyah/simple-pos/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
	"gorm.io/gorm"
)

type AuthController struct {
	UserRepo    *repositories.UserRepository
	TokenRepo   *repositories.PersonalAccessTokenRepository
	AuthService service.AuthService
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		TokenRepo:   repositories.NewPersonalAccessTokenRepository(db),
		UserRepo:    repositories.NewUserRepository(db),
		AuthService: service.NewAuthService(db),
	}
}

// Login godoc
// @Summary Authenticate a user
// @Description Login with email and password to receive a personal access token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response{data=dto.LoginResponse}
// @Failure 401 {object} utils.UnauthorizedResponse
// @Failure 422 {object} utils.ValidationErrorResponse
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	return ctrl.AuthService.Login(c)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and revoke current access token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SimpleResponse
// @Failure 401 {object} utils.UnauthorizedResponse
// @Router /auth/logout [post]
func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	token := c.Locals("token").(models.PersonalAccessToken)
	ctrl.TokenRepo.Delete(&token)

	return utils.SimpleSuccessResponse(c, "Logout successful. Current access token has been revoked.")
}

// ValidateToken godoc
// @Summary Validate authentication token
// @Description Validate the personal access token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.ValidateTokenResponse}
// @Failure 401 {object} utils.UnauthorizedResponse
// @Router /auth/validate-token [get]
func (ctrl *AuthController) ValidateToken(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)
	user, err := ctrl.UserRepo.FindByID(userId)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User not found")
	}

	return utils.SuccessResponse(c, "Token is valid", dto.ValidateTokenResponse{
		User: dto.ValidateTokenUserResponse{
			ID:   user.ID,
			Code: user.Code,
			Name: user.Name,
		},
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user password with current password and new password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param change-password body dto.ChangePasswordRequest true "Change password"
// @Success 200 {object} utils.Response{data=dto.LoginResponse}
// @Failure 401 {object} utils.UnauthorizedResponse
// @Failure 422 {object} utils.ValidationErrorResponse
// @Router /auth/change-password [post]
func (ctrl *AuthController) ChangePassword(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	user, err := ctrl.UserRepo.FindByID(userId)

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

	newToken, err := ctrl.AuthService.ChangePassword(
		user,
		data["current_password"].(string),
		data["new_password"].(string),
	)

	if err != nil {
		if err.Error() == "current_password_incorrect" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Your current password is incorrect")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update password")
	}

	return utils.SuccessResponse(c, "Password changed successfully", dto.LoginResponse{
		Token: newToken,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile with name and email
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body dto.UpdateProfileRequest true "Update profile"
// @Success 200 {object} utils.SimpleResponse
// @Failure 400 {object} utils.SimpleErrorResponse
// @Failure 401 {object} utils.UnauthorizedResponse
// @Failure 422 {object} utils.ValidationErrorResponse
// @Router /auth/profile [put]
func (ctrl *AuthController) UpdateProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	user, err := ctrl.UserRepo.FindByID(userId)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: User not found")
	}

	var req dto.UpdateProfileRequest

	rules := govalidator.MapData{
		"name":  []string{"required", "min:3"},
		"email": []string{"email"},
	}

	errs := utils.ValidateJSON(c, &req, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
	}

	err = ctrl.AuthService.UpdateProfile(
		user,
		req.Name,
		req.Email,
	)

	if err != nil {
		if err.Error() == "email_already_used" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email already used")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile")
	}

	return utils.SimpleSuccessResponse(c, "Profile updated successfully")
}
