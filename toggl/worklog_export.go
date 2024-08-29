package toggl

import (
	apphttp "JiraWorklogsImporter/http"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ExportWorkLogs(apiToken string, userId string, clientId string, workspaceId string, since string, until string) ([][]string, error) {
	baseURL := "https://track.toggl.com/reports/api/v2/details.csv"

	// Prepare query parameters
	params := url.Values{}
	params.Add("rounding", "Off")
	params.Add("sortDirection", "asc")
	params.Add("sortBy", "date")
	params.Add("order_field", "date")
	params.Add("order_desc", "off")
	params.Add("user_ids", userId)
	params.Add("client_ids", clientId)
	params.Add("since", since)
	params.Add("until", until)
	params.Add("billable", "both")
	params.Add("workspace_id", workspaceId)
	params.Add("show_amounts", "no")
	params.Add("date_format", "MM/DD/YYYY")
	params.Add("duration_format", "improved")
	params.Add("display_hours", "improved")
	params.Add("status", "active")
	params.Add("calculate", "time")
	params.Add("datepage", "1")
	params.Add("subgrouping", "users")
	params.Add("distinct_rates", "Off")
	params.Add("period", "thisWeek")
	params.Add("with_total_currencies", "1")
	params.Add("user_agent", "Snowball")
	params.Add("bars_count", "31")
	params.Add("subgrouping_ids", "true")

	// Combine base URL with parameters
	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Perform the request
	client := &http.Client{}
	resp, err := apphttp.MakeRequest(client, "GET", finalURL, apiToken, "api_token", []byte{})
	if err != nil {
		return [][]string{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return [][]string{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return [][]string{}, fmt.Errorf("exporting failed with status code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(strings.NewReader(string(body)))
	reader.Comma = ','
	reader.LazyQuotes = true
	return reader.ReadAll()
}
