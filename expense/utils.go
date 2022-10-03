package expense

import (
	"strings"
	"time"
)

const (
	DATE_PATTERN = "02-01-2006"
)

func ParseExpenseDate(dateStr string) (*time.Time, error) {
	date := time.Now()
	if "today" == strings.ToLower(dateStr) {
		return &date, nil
	}

	date, err := time.Parse(DATE_PATTERN, dateStr)
	if err != nil {
		return nil, err
	}

	return &date, nil
}
