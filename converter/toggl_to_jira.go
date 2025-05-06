package converter

import (
	"os"
)

type TogglToJiraConverter struct{}

func (c *TogglToJiraConverter) Convert(record []string) (ConvertedRecord, error) {
	descriptionRegex, exists := os.LookupEnv("DESCRIPTION_REGEX")
	if !exists {
		descriptionRegex = `^(.*?)\s*(?:\((.*?)\))?$`
	}

	description := record[5]
	durationString := record[11]

	startedAtDateTime, err := ConvertDateFormat(record[7] + " " + record[8])
	if err != nil {
		return ConvertedRecord{}, err
	}

	issueIdOrKey, contentText, err := ConvertToIssueIdAndContextText(description, descriptionRegex)
	if err != nil {
		return ConvertedRecord{}, err
	}

	timeSpentSeconds, err := ConvertToSeconds(durationString)
	if err != nil {
		return ConvertedRecord{}, err
	}

	return ConvertedRecord{
		IssueIdOrKey:      issueIdOrKey,
		ContentText:       contentText,
		StartedAtDateTime: startedAtDateTime,
		TimeSpentSeconds:  timeSpentSeconds,
	}, nil
}

func (c *TogglToJiraConverter) Supports(strategy string) bool {
	return strategy == "toggl_to_jira" || strategy == "csv_to_jira"
}
