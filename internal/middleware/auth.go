package middleware

import (
	"novaardiansyah/simple-pos/internal/service"
	"novaardiansyah/simple-pos/pkg/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"novaardiansyah/simple-pos/internal/repositories"
)

func Auth(db *gorm.DB) fiber.Handler {
	PersonalAccessTokenRepo := repositories.NewPersonalAccessTokenRepository(db)
	authService := service.NewAuthService(db)

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: No token provided")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: Invalid token format")
		}

		token, _, err := authService.ValidateToken(tokenString, "auth_token")

		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}

		fields := map[string]interface{}{"last_used_at": time.Now()}
		PersonalAccessTokenRepo.UpdateFields(token, fields)

		UserId := token.TokenableID

		c.Locals("token", *token)
		c.Locals("user_id", UserId)

		return c.Next()
	}
}
