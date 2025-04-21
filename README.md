# retrient

A lightweight REST client in Go with built-in support for:
- Retry logic with exponential backoff
- Rate limiting (based on `golang.org/x/time/rate`)
- Request timeouts
- Custom headers

The main goal is to support automatic retries for transient network errors and HTTP 5xx responses.

### How to use:

To initialise the client, the following can be used:

```go
c := client.New(
    "google.com",                            // baseURL to send request to 
    rate.Every(500*time.Millisecond),        // 2 requests per second
    2,                                       // burst
    10*time.Second,                          // timeout
    3,                                       // maximum retries
    1*time.Second,                           // backoff interval
    map[string]string{
        "Content-Type": "application/json",
    },
)
```

To send the request:
```go
reqMethod := "GET" 
path := "path-of-request"
resp, err := c.DoRequest(context.Background(), reqMethod, path, nil)
if err != nil {
    log.Fatalf("Request failed: %v", err)
}
defer resp.Body.Close()

body, _ := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

The `main.go` file contains this example as well for easier testing.