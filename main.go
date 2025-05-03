package main

import (
	"JiraWorklogsImporter/clockify"
	"JiraWorklogsImporter/converter"
	"JiraWorklogsImporter/importer"
	"JiraWorklogsImporter/jira"
	"JiraWorklogsImporter/toggl"
	"bufio"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	var project string
	var records [][]string
	var csvFilePathToImport string
	var since string
	var until string
	var nonInteractive bool

	flag.StringVar(&project, "project", "", "A project name that will load .env.<project-name> file.")
	flag.StringVar(&csvFilePathToImport, "import", "", "CSV file path to import.")
	flag.StringVar(&since, "since", "", "Import work logs since date. Format YYYY-MM-DD.")
	flag.StringVar(&until, "until", "", "Import work logs until date. Format YYYY-MM-DD.")
	flag.BoolVar(&nonInteractive, "non-interactive", false, "Non-interactive mode.")
	flag.BoolVar(&nonInteractive, "n", false, "An alias of --non-interactive.")
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file.")
		return
	}
	if project != "" {
		projectEnv := fmt.Sprintf(".env.%s", project)
		err = godotenv.Load(projectEnv)
		if err != nil {
			fmt.Printf("Error loading %s file.\n", projectEnv)
			return
		}
	}

	optionsValidationFailed := false

	if since == "" {
		fmt.Println("Missing since option.")
		optionsValidationFailed = true
	} else if !checkDateFormat(since) {
		fmt.Println("Invalid since option. The date must be in YYYY-MM-DD format.")
		optionsValidationFailed = true
	}

	if until == "" {
		fmt.Println("Missing until option.")
		optionsValidationFailed = true
	} else if !checkDateFormat(until) {
		fmt.Println("Invalid until option. The date must be in YYYY-MM-DD format.")
		optionsValidationFailed = true
	}

	if optionsValidationFailed {
		return
	}

	atlassianDomain := os.Getenv("ATLASSIAN_DOMAIN")
	atlassianEmail := os.Getenv("ATLASSIAN_EMAIL")
	atlassianApiToken := os.Getenv("ATLASSIAN_API_TOKEN")

	importStrategy, exists := os.LookupEnv("IMPORT_STRATEGY")
	if !exists {
		importStrategy = "csv_to_jira"
	}

	if importStrategy == "csv_to_jira" {
		if csvFilePathToImport == "" {
			fmt.Println("The CSV file is not provided. Use --import option.", err)
			return
		}

		records, err = importer.ReadCSVFile(csvFilePathToImport)
		if err != nil {
			fmt.Println("Error reading CSV file:", err)
			return
		}
	} else if importStrategy == "toggl_to_jira" {
		records, err = toggl.ExportWorkLogs(
			os.Getenv("TOGGL_API_TOKEN"),
			os.Getenv("TOGGL_USER_ID"),
			os.Getenv("TOGGL_CLIENT_ID"),
			os.Getenv("TOGGL_WORKSPACE_ID"),
			since,
			until,
		)
	} else if importStrategy == "clockify_to_jira" {
		records, err = clockify.ExportWorkLogs(
			os.Getenv("CLOCKIFY_API_TOKEN"),
			os.Getenv("CLOCKIFY_USER_ID"),
			os.Getenv("CLOCKIFY_PROJECT_ID"),
			os.Getenv("CLOCKIFY_WORKSPACE_ID"),
			since,
			until,
		)
	} else {
		fmt.Println("The given import strategy is not supported.")
		return
	}

	tableWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	for _, record := range records {
		fmt.Fprintln(tableWriter, strings.Join(record, "\t"))
	}

	tableWriter.Flush()

	if nonInteractive == false {
		fmt.Print("Please confirm the import [y/N]: ")
		reader := bufio.NewReader(os.Stdin)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred while reading the input. Please try again.", err)
			return
		}

		input = strings.TrimSpace(strings.ToLower(input))
		confirmed := input == "y"

		if confirmed == false {
			return
		}
	}

	for recordNo, record := range records {
		// Skip headers
		if recordNo == 0 {
			continue
		}

		factory := converter.NewConverterFactory()
		supportedConverter, err := factory.GetConverter(importStrategy)
		if err != nil {
			fmt.Println(err)
			continue
		}
		convertedRecord, err := supportedConverter.Convert(record)
		if err != nil {
			fmt.Println(err)
			continue
		}

		jira.ImportWorkLog(
			atlassianDomain,
			atlassianEmail,
			atlassianApiToken,
			convertedRecord.IssueIdOrKey,
			convertedRecord.ContentText,
			convertedRecord.StartedAtDateTime,
			convertedRecord.TimeSpentSeconds,
			recordNo,
		)
	}

	tableWriter.Flush()
}

func checkDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)

	return err == nil
}
