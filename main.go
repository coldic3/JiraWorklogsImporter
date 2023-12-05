package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

	reader := bufio.NewReader(os.Stdin)

	// Fetching default values from environment variables
	domain := os.Getenv("ATLASSIAN_DOMAIN")
	email := os.Getenv("EMAIL")
	apiToken := os.Getenv("API_TOKEN")
	defaultIssueIdOrKey := os.Getenv("DEFAULT_ISSUE_ID_OR_KEY")
	defaultContentText := os.Getenv("DEFAULT_CONTENT_TEXT")
	defaultDateString := os.Getenv("DEFAULT_DATE_STRING")
	defaultDurationString := os.Getenv("DEFAULT_DURATION_STRING")

	// Prompt for input with defaults
	issueIdOrKey := promptWithDefault(reader, "Enter issue ID or key: ", defaultIssueIdOrKey)
	contentText := promptWithDefault(reader, "Enter content text: ", defaultContentText)
	dateString := promptWithDefault(reader, "Enter date (DD.MM.YYYY): ", defaultDateString)
	durationString := promptWithDefault(reader, "Enter time duration (e.g., 1h 30m, 45m, 8h): ", defaultDurationString)

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
	payload.Started = started
	payload.TimeSpentSeconds = timeSpentSeconds

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

	fmt.Println("\033[1;32mTIME LOGGED!\033[0m")
}

func promptWithDefault(reader *bufio.Reader, prompt string, defaultValue string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
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
