package clockify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TimeEntry struct {
	Description  string       `json:"description"`
	TimeInterval TimeInterval `json:"timeInterval"`
}

type TimeInterval struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func ExportWorkLogs(apiToken string, userId string, projectId string, workspaceId string, since string, until string) ([][]string, error) {
	baseURL := fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/user/%s/time-entries", workspaceId, userId)

	// Prepare query parameters
	url := fmt.Sprintf("%s?project=%s&start=%sT00:00:00Z&end=%sT23:59:59Z", baseURL, projectId, since, until)

	// Perform the request
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiToken)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return [][]string{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return [][]string{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return [][]string{}, fmt.Errorf("exporting failed with status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var entries []TimeEntry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		return [][]string{}, err
	}

	// Prepare CSV-like output
	result := [][]string{{"Description", "Start", "End"}}
	for _, entry := range entries {
		result = append(result, []string{
			entry.Description,
			entry.TimeInterval.Start.Format("2006-01-02 15:04:05"),
			entry.TimeInterval.End.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}
