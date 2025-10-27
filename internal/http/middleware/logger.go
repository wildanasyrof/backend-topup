package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/wildanasyrof/backend-topup/pkg/logger" // Import your logger package
)

func LoggerMiddleware(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		started := time.Now()

		// ---- Correlation / Request ID ---------------------------------------
		reqID := c.Get(fiber.HeaderXRequestID)
		if reqID == "" {
			if v := c.Locals("requestid"); v != nil {
				if s, ok := v.(string); ok {
					reqID = s
				}
			}
		}
		if reqID == "" {
			reqID = utils.UUID() // light dep from fiber/utils
		}
		c.Locals("requestid", reqID)
		c.Set(fiber.HeaderXRequestID, reqID)

		// ---- BEST PRACTICE: Create and inject request-scoped logger ----
		reqLogger := log.With(logger.Fields{
			"request_id": reqID,
		})
		c.Locals(logger.CtxKey, reqLogger) // Use the const from logger package

		// ---- Next handler ----------------------------------------------------
		err := c.Next()

		// ---- Request/Response facts -----------------------------------------
		latency := time.Since(started)
		method := c.Method()
		path := c.OriginalURL()
		status := c.Response().StatusCode()

		routePath := ""
		if r := c.Route(); r != nil {
			routePath = r.Path // e.g. /v1/deposits/:id
		}

		ip := c.IP()
		realIP := ""
		if ips := c.IPs(); len(ips) > 0 {
			realIP = ips[0]
		}

		ua := string(c.Request().Header.UserAgent())
		referer := c.Get(fiber.HeaderReferer)
		host := c.Hostname()
		proto := c.Protocol()

		size := c.Response().Header.ContentLength()
		if size < 0 {
			size = len(c.Response().Body())
		}

		uid, _ := c.Locals("user_id").(string) // set by your auth middleware

		// ---- BEST PRACTICE: Log as structured fields, not one string ----
		fields := logger.Fields{
			"method":     method,
			"path":       path,
			"route":      routePath,
			"status":     status,
			"latency":    latency, // zerolog handles time.Duration
			"ip":         ip,
			"real_ip":    realIP,
			"host":       host,
			"protocol":   proto,
			"size_bytes": size,
			"user_id":    uid,
			"user_agent": ua,
			"referer":    referer,
		}

		// Simple, human-readable message
		msg := fmt.Sprintf("HTTP %s %s", method, path)

		// Log with correct level
		if status >= 500 {
			// Use the handler error if it exists, otherwise create one
			logErr := err
			if logErr == nil {
				logErr = fmt.Errorf("server error: status %d", status)
			}
			reqLogger.With(fields).Error(logErr, msg)
		} else if status >= 400 {
			reqLogger.With(fields).Warn(msg)
		} else {
			reqLogger.With(fields).Info(msg)
		}

		return err // Pass the original error up
	}
}
