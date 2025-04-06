package converter

import (
	"JiraWorklogsImporter/toggl"
	"fmt"
	"os"
)

type TogglToJiraConverter struct{}

func (c *TogglToJiraConverter) Convert(record []string) (string, error) {
	descriptionRegex, exists := os.LookupEnv("DESCRIPTION_REGEX")
	if !exists {
		descriptionRegex = `^(.*?)\s*(?:\((.*?)\))?$`
	}

	description := record[5]
	durationString := record[11]

	startedAtDateTime, err := toggl.ConvertDateFormat(record[7] + " " + record[8])
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

func (c *TogglToJiraConverter) Supports(format string) bool {
	return format == "toggl_to_jira"
}
