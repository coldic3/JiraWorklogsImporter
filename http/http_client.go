package http

import (
	"bytes"
	"net/http"
)

func MakeRequest(client *http.Client, method string, url string, user string, password string, jsonPayload []byte) (*http.Response, error) {
	// Create a new request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Set basic auth
	req.SetBasicAuth(user, password)

	// Perform the request
	return client.Do(req)
}
