package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateTopupID creates a unique ID for a top-up transaction.
// The format is "TOPUP" + timestamp + a random number.
func GenerateTopupID() string {
	// Use a fixed prefix for clarity.
	prefix := "TOPUP"

	// Get the current timestamp in a specific format (e.g., YYYYMMDDHHmmss).
	timestamp := time.Now().Format("01021504")

	// Generate a small random number to prevent collisions if requests happen at the same millisecond.
	randomNumber := rand.Intn(10000)

	// Combine all parts into a single string.
	return fmt.Sprintf("%s%s%d", prefix, timestamp, randomNumber)
}
