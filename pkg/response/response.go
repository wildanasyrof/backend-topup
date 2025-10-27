// pkg/response/response.go
package response

import "github.com/gofiber/fiber/v2"

type Envelope struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	Error     *ErrBody    `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func OK(c *fiber.Ctx, data interface{}, meta ...interface{}) error {
	env := Envelope{Success: true, Data: data, RequestID: requestID(c)}
	if len(meta) > 0 && meta[0] != nil {
		env.Meta = meta[0]
	}
	return c.Status(fiber.StatusOK).JSON(env)
}

func Created(c *fiber.Ctx, data interface{}, meta ...interface{}) error {
	env := Envelope{Success: true, Data: data, RequestID: requestID(c)}
	if len(meta) > 0 && meta[0] != nil {
		env.Meta = meta[0]
	}
	return c.Status(fiber.StatusCreated).JSON(env)
}

func requestID(c *fiber.Ctx) string {
	// Works with fiber/middleware/requestid or any upstream proxy setting the header
	if id := c.Get("X-Request-ID"); id != "" {
		return id
	}
	return c.Locals("requestid").(string) // fiber requestid middleware sets this
}
