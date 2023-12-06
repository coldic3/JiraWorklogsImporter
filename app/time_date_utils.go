package app

import (
	"fmt"
	"time"
)

func FormatDate(dateStr string) (string, error) {
	parsedDate, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02T08:00:00.000+0000"), nil
}

func ParseDuration(durationStr string) (int, error) {
	var hours, minutes int
	_, err := fmt.Sscanf(durationStr, "%dh %dm", &hours, &minutes)
	if err != nil {
		_, err := fmt.Sscanf(durationStr, "%dh", &hours)
		if err != nil {
			_, err := fmt.Sscanf(durationStr, "%dm", &minutes)
			if err != nil {
				return 0, err
			}
		}
	}
	return hours*3600 + minutes*60, nil
}
