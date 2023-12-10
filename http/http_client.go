package http

import (
	"bytes"
	"net/http"
)

func MakeRequest(client *http.Client, url string, email string, apiToken string, jsonPayload []byte) (*http.Response, error) {
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
