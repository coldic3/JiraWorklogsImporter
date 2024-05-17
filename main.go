package main

import (
	"JiraWorklogsImporter/importer"
	"JiraWorklogsImporter/jira"
	"JiraWorklogsImporter/toggl"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"text/tabwriter"
)

func main() {
	var records [][]string
	var csvFilePathToImport string
	var since string
	var until string
	var dryRun bool

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file.")
		return
	}

	flag.StringVar(&csvFilePathToImport, "import", "", "CSV file path to import.")
	flag.StringVar(&since, "since", "", "Import work logs since date. Format YYYY-MM-DD.")
	flag.StringVar(&until, "until", "", "Import work logs until date. Format YYYY-MM-DD.")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run. Export work logs but do not import.")
	flag.Parse()

	optionsValidationFailed := false

	if since == "" {
		fmt.Println("Missing since option.")
		optionsValidationFailed = true
	} else if !toggl.CheckDateFormat(since) {
		fmt.Println("Invalid since option. The date must be in YYYY-MM-DD format.")
		optionsValidationFailed = true
	}

	if until == "" {
		fmt.Println("Missing until option.")
		optionsValidationFailed = true
	} else if !toggl.CheckDateFormat(until) {
		fmt.Println("Invalid until option. The date must be in YYYY-MM-DD format.")
		optionsValidationFailed = true
	}

	if optionsValidationFailed {
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

	tableWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	for recordNo, record := range records {
		if dryRun {
			fmt.Fprintln(tableWriter, strings.Join(record, "\t"))
			continue
		}

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

	tableWriter.Flush()
}
