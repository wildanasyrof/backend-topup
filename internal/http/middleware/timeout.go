package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
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
			return apperror.New(apperror.CodeTimeout, "REQUEST_TIME_OUT", nil)
		case err := <-errCh:
			return err
		}
	}
}
