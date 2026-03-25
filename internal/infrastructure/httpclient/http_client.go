package httpclient

import (
	"net/http"
	"time"
)

func NewHttpClient(transport http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: transport,
		Timeout: 1 * time.Minute,
	}
}