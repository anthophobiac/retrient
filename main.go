package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/time/rate"
	"retrient/client"
)

var baseURL = "https://jsonplaceholder.typicodme.com"
var path = "/posts/1"

func main() {
	c := client.New(
		baseURL,
		rate.Every(500*time.Millisecond),
		2,
		10*time.Second,
		3,
		1*time.Second,
		map[string]string{
			"Content-Type": "application/json",
		},
	)

	resp, err := c.DoRequest(context.Background(), "GET", path, nil)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Body:", string(body))
}
