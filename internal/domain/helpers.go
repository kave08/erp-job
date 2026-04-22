package domain

import (
	"fmt"
	"strconv"
)

func ParseVisitorCode(raw string) (int, error) {
	visitorCode, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse visitor code %q: %w", raw, err)
	}

	return visitorCode, nil
}
