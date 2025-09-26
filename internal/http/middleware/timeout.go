package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/pkg/response"
)

func TimeoutMiddleware(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), d)
		defer cancel()

		// replace Fiber's context with ours
		c.SetUserContext(ctx)

		errCh := make(chan error, 1)
		go func() {
			errCh <- c.Next()
		}()

		select {
		case <-ctx.Done():
			return response.Error(c, fiber.StatusGatewayTimeout, "request time out", nil)
		case err := <-errCh:
			return err
		}
	}
}
