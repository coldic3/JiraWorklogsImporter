package main

import (
	"JiraWorklogsImporter/importer"
	"JiraWorklogsImporter/jira"
	"JiraWorklogsImporter/toggl"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	var records [][]string
	var csvFilePathToImport string
	var since string
	var until string

	flag.StringVar(&csvFilePathToImport, "import", "", "CSV file path to import")
	flag.StringVar(&since, "since", "", "Import work logs since date")
	flag.StringVar(&until, "until", "", "Import work logs until date")
	flag.Parse()

	if since == "" {
		fmt.Println("Missing since option")
		return
	}

	if until == "" {
		fmt.Println("Missing until option")
		return
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	atlassianDomain := os.Getenv("ATLASSIAN_DOMAIN")
	atlassianEmail := os.Getenv("ATLASSIAN_EMAIL")
	atlassianApiToken := os.Getenv("ATLASSIAN_API_TOKEN")
	togglApiToken := os.Getenv("TOGGL_API_TOKEN")
	togglUserId := os.Getenv("TOGGL_USER_ID")
	togglClientId := os.Getenv("TOGGL_CLIENT_ID")
	togglWorkspaceId := os.Getenv("TOGGL_WORKSPACE_ID")

	if csvFilePathToImport != "" {
		records, err = importer.ReadCSVFile(csvFilePathToImport)
		if err != nil {
			fmt.Println("Error reading CSV file:", err)
			return
		}
	} else {
		records, err = toggl.ExportWorkLogs(togglApiToken, togglUserId, togglClientId, togglWorkspaceId, since, until)
	}

	for recordNo, record := range records {
		// Skip headers
		if recordNo == 0 {
			continue
		}

		description := record[5]
		durationString := record[11]

		startedAtDateTime, err := toggl.ConvertDateFormat(record[7] + " " + record[8])
		if err != nil {
			fmt.Println(fmt.Sprintln(err))
			continue
		}

		issueIdOrKey, contentText, err := toggl.ConvertToIssueIdAndContextText(description)
		if err != nil {
			fmt.Println(fmt.Sprintln(err))
			continue
		}

		timeSpentSeconds, err := toggl.ConvertToSeconds(durationString)
		if err != nil {
			fmt.Println(fmt.Sprintln(err))
			continue
		}

		jira.ImportWorkLog(atlassianDomain, atlassianEmail, atlassianApiToken, issueIdOrKey, contentText, startedAtDateTime, timeSpentSeconds, recordNo)
	}
}
