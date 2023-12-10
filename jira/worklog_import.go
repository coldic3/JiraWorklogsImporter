package jira

import (
	apphttp "JiraWorklogsImporter/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ImportWorkLog(domain string, email string, apiToken string, issueIdOrKey string, contentText string, date string, time int, recordNo int) {
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
	resp, err := apphttp.MakeRequest(client, url, email, apiToken, jsonPayload)
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
		fmt.Printf("\033[1;31mERROR!\u001B[0m Record no %d has not been imported.\n", recordNo)
	}
}
