package helpers

import (
	"strconv"

	errors "entain/internal/error"
)

// ExtractInt64Route - returns the route int64 variable for the current request.
func ExtractInt64Route(vars map[string]string, key string) (uint64, error) {
	str, ok := vars[key]
	if !ok {
		return 0, errors.ErrBadRouting
	}
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint64(result), nil
}
