package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

type Payload struct {
	Comment struct {
		Content []struct {
			Content []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"content"`
			Type string `json:"type"`
		} `json:"content"`
		Type    string `json:"type"`
		Version int    `json:"version"`
	} `json:"comment"`
	Started          string `json:"started"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

func makeRequest(client *http.Client, url string, email string, apiToken string, jsonPayload []byte) (*http.Response, error) {
	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Set basic auth
	req.SetBasicAuth(email, apiToken)

	// Perform the request
	return client.Do(req)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Fetching default values from environment variables
	domain := os.Getenv("ATLASSIAN_DOMAIN")
	email := os.Getenv("EMAIL")
	apiToken := os.Getenv("API_TOKEN")

	csvFilePath := "import_me.csv"
	records, err := readCSVFile(csvFilePath)
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
		started, err := formatDate(dateString)
		if err != nil {
			fmt.Println("Error in date formatting:", err)
			return
		}

		// Convert duration to seconds
		timeSpentSeconds, err := parseDuration(durationString)
		if err != nil {
			fmt.Println("Error in parsing duration:", err)
			return
		}

		importWorkLog(domain, email, apiToken, issueIdOrKey, contentText, started, timeSpentSeconds, recordNo)
	}
}

func importWorkLog(domain string, email string, apiToken string, issueIdOrKey string, contentText string, date string, time int, recordNo int) {
	// Create the URL
	url := fmt.Sprintf("https://%s.atlassian.net/rest/api/3/issue/%s/worklog", domain, issueIdOrKey)

	// Set up the payload
	payload := Payload{}
	payload.Comment.Content = []struct {
		Content []struct {
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"content"`
		Type string `json:"type"`
	}{
		{
			Content: []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			}{
				{
					Text: contentText,
					Type: "text",
				},
			},
			Type: "paragraph",
		},
	}
	payload.Comment.Type = "doc"
	payload.Comment.Version = 1
	payload.Started = date
	payload.TimeSpentSeconds = time

	jsonPayload, _ := json.Marshal(payload)

	// Perform the request
	client := &http.Client{}
	resp, err := makeRequest(client, url, email, apiToken, jsonPayload)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var prettyJson bytes.Buffer
	prettyJsonError := json.Indent(&prettyJson, body, "", "    ")
	if prettyJsonError != nil {
		fmt.Println("JSON parse prettyJsonError: ", prettyJsonError)
		return
	}

	if resp.StatusCode == 201 {
		fmt.Println("\033[1;32mTIME LOGGED!\033[0m")
	} else {
		fmt.Printf("\033[1;31mERROR!\u001B[0m Record no %d has not been imported.\n", recordNo+1)
	}
}

func readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	return reader.ReadAll()
}

func formatDate(dateStr string) (string, error) {
	parsedDate, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02T08:00:00.000+0000"), nil
}

func parseDuration(durationStr string) (int, error) {
	var hours, minutes int
	_, err := fmt.Sscanf(durationStr, "%dh %dm", &hours, &minutes)
	if err != nil {
		_, err := fmt.Sscanf(durationStr, "%dh", &hours)
		if err != nil {
			_, err := fmt.Sscanf(durationStr, "%dm", &minutes)
			if err != nil {
				return 0, err
			}
		}
	}
	return hours*3600 + minutes*60, nil
}
