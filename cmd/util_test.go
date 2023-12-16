package cmd

import (
	"net/http"

	"github.com/dstotijn/go-notion"
)

type MockRoundTripper struct {
	Response *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, nil
}

func NewMockClient(mockResponse *http.Response, apiKey string) *notion.Client {
	return notion.NewClient(apiKey, notion.WithHTTPClient(&http.Client{
		Transport: &MockRoundTripper{Response: mockResponse},
	},
	))
}
