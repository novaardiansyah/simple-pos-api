/*
 * Project Name: controllers
 * File: user_controller.go
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
	"novaardiansyah/simple-pos/internal/repositories"
	"novaardiansyah/simple-pos/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	repo *repositories.UserRepository
}

func NewUserController(repo *repositories.UserRepository) *UserController {
	return &UserController{repo: repo}
}

// Index godoc
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(15)
// @Success 200 {object} utils.PaginatedResponse{data=[]UserSwagger}
// @Failure 500 {object} utils.Response
// @Router /users [get]
// @Security BearerAuth
func (ctrl *UserController) Index(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	total, err := ctrl.repo.Count()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count users")
	}

	users, err := ctrl.repo.FindAllPaginated(page, perPage)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users")
	}

	return utils.PaginatedSuccessResponse(c, "Users retrieved successfully", users, page, perPage, total, len(users))
}

// Show godoc
// @Summary Get user details
// @Description Get detailed information about a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{data=UserSwagger}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [get]
// @Security BearerAuth
func (ctrl *UserController) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
}
