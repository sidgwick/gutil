# Golang HTTP util

## Example usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/sidgwick/gutil/http"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, TRACE_KEY, "this-is-trace-id")

	initLogger() // see project logger init example

	resp, err := runHttpRequest(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", resp)
}

func runHttpRequest(ctx context.Context) (interface{}, error) {
	var resp interface{}

	client, err := http.NewCurl("https://api.seeip.org/jsonip?")
	if err != nil {
		return nil, err
	}

	req := &http.Request{}
	err = client.GetJson(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
```