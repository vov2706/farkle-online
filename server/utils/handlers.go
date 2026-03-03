package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func ConvertRouteParamToUint(c fiber.Ctx, key string) (uint, error) {
	idStr := c.Params(key)
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return 1, fmt.Errorf("invalid route param %s", key)
	}

	return uint(id), nil
}

func ConvertQueryParamToUint(c fiber.Ctx, key string) (uint, error) {
	idStr := c.Query(key)
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return 1, fmt.Errorf("invalid query param %s", key)
	}

	return uint(id), nil
}
