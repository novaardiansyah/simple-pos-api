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
	return ctrl.AuthService.ChangePassword(c)
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

// RefreshToken godoc
// @Summary Refresh token
// @Description Refresh token with current token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param refresh-token body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.Response{data=dto.LoginResponse}
// @Failure 401 {object} utils.UnauthorizedResponse
// @Failure 422 {object} utils.ValidationErrorResponse
// @Router /auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *fiber.Ctx) error {
	return ctrl.AuthService.RefreshToken(c)
}
