package request

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ParseQueryToUUID returns pointer UUID or error from query request.
func ParseQueryToUUID(c *fiber.Ctx, key string) (*uuid.UUID, error) {
	const op = "lib.api.request.parseQueryToUUID"

	str := c.Query(key)

	if str == "" {
		return nil, nil //nolint:nilnil
	}

	v, err := uuid.Parse(str)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &v, nil
}

// ParseQueryToTime returns pointer to Time and error from query request by layout.
func ParseQueryToTime(c *fiber.Ctx, key string, layout string) (*time.Time, error) {
	const op = "lib.api.request.parseQueryToTime"

	str := c.Query(key)

	if str == "" {
		return nil, nil //nolint:nilnil
	}

	v, err := time.Parse(layout, str)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &v, nil
}
