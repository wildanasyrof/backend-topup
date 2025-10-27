package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/jwt"
)

func Auth(jwtSvc jwt.JWTService, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return apperror.ErrUnauthorized
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		id, role, err := jwtSvc.ValidateToken(tokenStr)
		if err != nil {
			return apperror.ErrUnauthorized
		}

		if role == "" {
			return apperror.ErrUnauthorized
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
		return apperror.ErrForbidden
	}
}
