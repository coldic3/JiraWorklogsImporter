package toggl

import (
	"fmt"
	"regexp"
	"time"
)

func ConvertToSeconds(hmsTime string) (int, error) {
	parsedTime, err := time.Parse("15:04:05", hmsTime)

	if err != nil {
		return 0, fmt.Errorf("cannot parse string \"%s\" into duration", hmsTime)
	}

	hours := parsedTime.Hour()
	minutes := parsedTime.Minute()
	seconds := parsedTime.Second()

	return hours*3600 + minutes*60 + seconds, nil
}

func ConvertDateFormat(date string) (string, error) {
	timeZone := "0100"

	parsedTime, err := time.Parse("2006-01-02 15:04:05", date)

	if err != nil {
		return "", fmt.Errorf("cannot parse string \"%s\" into date format yyyy-MM-dd'T'HH:mm:ss.SSSZ", date)
	}

	return parsedTime.Format("2006-01-02T15:04:05.000") + "+" + timeZone, nil
}

func ConvertToIssueIdAndContextText(text string) (string, string, error) {
	re := regexp.MustCompile(`^(.*?)\s*(?:\((.*?)\))?$`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid record format: %s", text)
	}

	return matches[1], matches[2], nil
}
