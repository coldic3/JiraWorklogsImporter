package main

import (
	"JiraWorklogsImporter/importer"
	"JiraWorklogsImporter/importer/toggl"
	"JiraWorklogsImporter/jira"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	var csvFilePathToImport string
	flag.StringVar(&csvFilePathToImport, "import", "", "CSV file path to import")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	domain := os.Getenv("ATLASSIAN_DOMAIN")
	email := os.Getenv("EMAIL")
	apiToken := os.Getenv("API_TOKEN")

	records, err := importer.ReadCSVFile(csvFilePathToImport)
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	for recordNo, record := range records {
		// Skip headers
		if recordNo == 0 {
			continue
		}

		description := record[5]
		durationString := record[11]
		startedAtDate := record[7]

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

		jira.ImportWorkLog(domain, email, apiToken, issueIdOrKey, contentText, startedAtDate, timeSpentSeconds, recordNo)
	}
}
