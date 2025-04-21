package client

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	BaseURL     string
	Headers     map[string]string
	HTTPClient  *http.Client
	RateLimiter *rate.Limiter
	MaxRetries  int
	Backoff     time.Duration
}

func New(
	baseURL string,
	rateLimit rate.Limit,
	burst int,
	timeout time.Duration,
	maxRetries int,
	backoff time.Duration,
	headers map[string]string,
) *Client {
	return &Client{
		BaseURL: baseURL,
		Headers: headers,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		RateLimiter: rate.NewLimiter(rateLimit, burst),
		MaxRetries:  maxRetries,
		Backoff:     backoff,
	}
}
