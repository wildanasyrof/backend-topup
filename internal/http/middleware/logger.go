// internal/http/middleware/request_logger.go
package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/wildanasyrof/backend-topup/pkg/logger"
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

		ip := c.IP()   // client IP (after proxy config)
		ips := c.IPs() // X-Forwarded-For chain
		realIP := ""
		if len(ips) > 0 {
			realIP = ips[0]
		}

		ua := string(c.Request().Header.UserAgent())
		referer := c.Get(fiber.HeaderReferer)
		host := c.Hostname()
		proto := c.Protocol() // "http" or "https"

		size := c.Response().Header.ContentLength()
		if size < 0 {
			size = len(c.Response().Body()) // fallback if no content-length
		}

		uid, _ := c.Locals("user_id").(string) // set by your auth middleware

		// Log handler error explicitly
		if err != nil {
			log.Error(err, "request error")
		}

		// ---- One-line, structured-ish summary -------------------------------
		summary := fmt.Sprintf(
			`%s %s -> %d in %s req_id=%s route=%s ip=%s real_ip=%s host=%s proto=%s size=%dB uid=%s ua=%q referer=%q`,
			method, path, status, latency.Truncate(time.Microsecond),
			reqID, routePath, ip, realIP, host, proto, size, uid, ua, referer,
		)

		// Elevate severity on 5xx
		if status >= 500 {
			if err == nil {
				log.Error(fmt.Errorf("http %d", status), summary)
			} else {
				log.Error(err, summary)
			}
		} else {
			log.Info(summary)
		}

		return err
	}
}
