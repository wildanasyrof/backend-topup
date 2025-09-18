package response

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, message string, data interface{}, codes ...int) error {
	code := fiber.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func Error(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"data":    data,
	})
}
