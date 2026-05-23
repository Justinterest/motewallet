package kun

import (
	"context"
	"strings"
	"time"
)

// WithRetry wraps a function with exponential backoff.
// Only retries on network errors and 5xx KUN errors. Never retries 4xx.
func WithRetry(ctx context.Context, maxRetries int, fn func() error) error {
	var lastErr error
	delay := time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if kunErr, ok := IsKUNError(lastErr); ok {
			if strings.HasPrefix(kunErr.Code, "4") {
				return lastErr
			}
		}

		if attempt == maxRetries {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		delay *= 2
		if delay > 10*time.Second {
			delay = 10 * time.Second
		}
	}

	return lastErr
}
