package service

import (
	"strings"
	"time"
)

func mapDepositOrderStatus(status string) string {
	switch strings.ToUpper(status) {
	case "SUCCESS":
		return "COMPLETED"
	case "FAIL", "FAILED":
		return "FAILED"
	case "PROCESSING", "PENDING":
		return "PROCESSING"
	default:
		return status
	}
}

func parseKunDateTime(value string) time.Time {
	if value == "" {
		return time.Now()
	}
	loc := kunDepositLocation()
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339} {
		if t, err := time.ParseInLocation(layout, value, loc); err == nil {
			return t
		}
	}
	return time.Now()
}

func kunDepositLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Local
	}
	return loc
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
