package converter

import (
	"JiraWorklogsImporter/toggl"
	"fmt"
	"os"
	"time"
)

type ClockifyToJiraConverter struct{}

func (c *ClockifyToJiraConverter) Convert(record []string) (string, error) {
	descriptionRegex, exists := os.LookupEnv("DESCRIPTION_REGEX")
	if !exists {
		descriptionRegex = `^(.*?)\s*(?:\((.*?)\))?$`
	}

	description := record[0]

	startTime, err := time.Parse("2006-01-02 15:04:05", record[1])
	if err != nil {
		return "", err
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", record[2])
	if err != nil {
		return "", err
	}

	duration := endTime.Sub(startTime)
	durationString := fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)

	startedAtDateTime, err := toggl.ConvertDateFormat(record[1])
	if err != nil {
		return "", err
	}

	issueIdOrKey, contentText, err := toggl.ConvertToIssueIdAndContextText(description, descriptionRegex)
	if err != nil {
		return "", err
	}

	timeSpentSeconds, err := toggl.ConvertToSeconds(durationString)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s,%s,%s,%d", issueIdOrKey, contentText, startedAtDateTime, timeSpentSeconds), nil
}

func (c *ClockifyToJiraConverter) Supports(format string) bool {
	return format == "clockify_to_jira"
}
