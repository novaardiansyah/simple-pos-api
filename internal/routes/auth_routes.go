package routes

import (
	"novaardiansyah/simple-pos/internal/controllers"
	"novaardiansyah/simple-pos/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoutes(api fiber.Router, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	auth := api.Group("/auth")
	auth.Use(middleware.AuthLimiter())

	auth.Post("/login", authController.Login)
	auth.Get("/validate-token", middleware.Auth(db), authController.ValidateToken)
	auth.Post("/logout", middleware.Auth(db), authController.Logout)
	auth.Post("/change-password", middleware.Auth(db), authController.ChangePassword)
	auth.Put("/profile", middleware.Auth(db), authController.UpdateProfile)
	auth.Post("/refresh", authController.RefreshToken)
}
