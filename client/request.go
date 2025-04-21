package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"time"

	"net/http"
)

func (c *Client) DoRequest(
	ctx context.Context,
	method string,
	path string,
	body []byte,
) (*http.Response, error) {
	url := c.BaseURL + path
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if err = c.RateLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		req, reqErr := http.NewRequestWithContext(
			ctx,
			method,
			url,
			bytes.NewBuffer(body),
		)
		if reqErr != nil {
			return nil, reqErr
		}

		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}

		resp, err = c.HTTPClient.Do(req)

		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		log.Printf("Attempt %d failed: %v", attempt+1, err)

		if resp != nil {
			_, err := io.Copy(io.Discard, resp.Body)
			if err != nil {
				return nil, err
			}
			_ = resp.Body.Close()
		}

		if attempt == c.MaxRetries {
			break
		}

		sleep := c.Backoff * (1 << attempt)
		time.Sleep(sleep)
	}

	if err != nil {
		return nil, err
	}
	return resp, errors.New("request failed with 5xx after retries")
}
