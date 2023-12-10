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

	return int(parsedTime.Sub(time.Time{}).Seconds()), nil
}

func ConvertToIssueIdAndContextText(text string) (string, string, error) {
	re := regexp.MustCompile(`^(.*?)\s*(?:\((.*?)\))?$`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid record format: %s", text)
	}

	return matches[1], matches[2], nil
}
