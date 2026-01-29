package routes

import (
	"novaardiansyah/simple-pos/internal/controllers"
	"novaardiansyah/simple-pos/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoutes(api fiber.Router, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	users := api.Group("/users", middleware.Auth(db))
	users.Get("/", userController.Index)
	users.Get("/me", userController.Me)
	users.Get("/:id", userController.Show)
}
