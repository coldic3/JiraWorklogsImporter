package main

import (
	"JiraWorklogsImporter/app"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

func main() {
	var csvFilePathToImport string
	flag.StringVar(&csvFilePathToImport, "file-path", "", "CSV file path to import")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Fetching default values from environment variables
	domain := os.Getenv("ATLASSIAN_DOMAIN")
	email := os.Getenv("EMAIL")
	apiToken := os.Getenv("API_TOKEN")

	records, err := app.ReadCSVFile(csvFilePathToImport)
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	for recordNo, record := range records {
		description := record[0]
		durationString := record[1]
		dateString := record[2]

		// Extract issueIdOrKey and contentText from description
		re := regexp.MustCompile(`^(.*?)\s*(?:\((.*?)\))?$`)
		matches := re.FindStringSubmatch(description)
		if len(matches) < 3 {
			fmt.Println("Invalid record format:", description)
			continue
		}
		issueIdOrKey := matches[1]
		contentText := matches[2]

		// Convert date to the required format
		started, err := app.FormatDate(dateString)
		if err != nil {
			fmt.Println("Error in date formatting:", err)
			return
		}

		// Convert duration to seconds
		timeSpentSeconds, err := app.ParseDuration(durationString)
		if err != nil {
			fmt.Println("Error in parsing duration:", err)
			return
		}

		app.ImportWorkLog(domain, email, apiToken, issueIdOrKey, contentText, started, timeSpentSeconds, recordNo)
	}
}
