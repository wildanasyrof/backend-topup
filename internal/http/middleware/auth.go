package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/pkg/jwt"
	"github.com/wildanasyrof/backend-topup/pkg/response"
)

func Auth(jwtSvc jwt.JWTService, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return response.Error(c, fiber.StatusUnauthorized, "Missing bearer token", nil)
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		id, role, err := jwtSvc.ValidateToken(tokenStr) // returns user id as string (e.g., JWT sub)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "invalid or expired refresh token", nil)
		}

		if role == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Role claim is missing in token", nil)
		}

		c.Locals("user_id", id)
		c.Locals("role", role)

		// If no specific roles are required, allow all authenticated users
		if len(allowedRoles) == 0 {
			return c.Next()
		}

		// Check if user has an allowed role
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next() // Proceed to the next handler
			}
		}

		// If role is not allowed, return forbidden
		return response.Error(c, fiber.StatusForbidden, "Unauthorized", "Forbidden")
	}
}
