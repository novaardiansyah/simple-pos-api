package routes

import (
	"novaardiansyah/simple-pos/internal/controllers"
	"novaardiansyah/simple-pos/internal/middleware"
	"novaardiansyah/simple-pos/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoutes(api fiber.Router, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo)

	users := api.Group("/users", middleware.Auth(db))
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)
}
