package connector

import (
	"fmt"
	"strconv"
)

func parsePageToken(token string) (uint, error) {
	if token == "" {
		return 0, nil
	}
	num, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("baton-oracle-scm: failed to parse page token: %w", err)
	}
	return uint(num), nil
}

func formatNextPageToken(offset int) string {
	if offset < 0 {
		return ""
	}
	return strconv.Itoa(offset)
}
