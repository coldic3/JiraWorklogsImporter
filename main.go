package main

import (
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
	descriptionRegex, exists := os.LookupEnv("DESCRIPTION_REGEX")
	if !exists {
		descriptionRegex = `^(.*?)\s*(?:\((.*?)\))?$`
	}

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

		description := record[5]
		durationString := record[11]

		startedAtDateTime, err := toggl.ConvertDateFormat(record[7] + " " + record[8])
		if err != nil {
			fmt.Println(fmt.Sprintln(err))
			continue
		}

		issueIdOrKey, contentText, err := toggl.ConvertToIssueIdAndContextText(description, descriptionRegex)
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
