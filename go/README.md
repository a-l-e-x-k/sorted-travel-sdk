# sorted-travel Go SDK

Official Go client and CLI for the [Sorted Travel](https://sorted.travel) REST API.

- Homepage: [https://sorted.travel](https://sorted.travel)
- Developer portal: [https://sorted.travel/developers](https://sorted.travel/developers)
- SDK guide: [https://sorted.travel/sdks.md](https://sorted.travel/sdks.md)
- Repository: [https://github.com/a-l-e-x-k/sorted-travel-sdk](https://github.com/a-l-e-x-k/sorted-travel-sdk)

## Install

```sh
go get github.com/a-l-e-x-k/sorted-travel-sdk/go/sortedtravel@v0.1.1
```

CLI:

```sh
go install github.com/a-l-e-x-k/sorted-travel-sdk/go/cmd/sorted-travel@v0.1.1
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"

	"github.com/a-l-e-x-k/sorted-travel-sdk/go/sortedtravel"
)

func main() {
	client := sortedtravel.NewClient()
	status, err := client.Status(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(status)

	page, err := client.ListDestinations(context.Background(), nil, intPtr(20))
	if err != nil {
		panic(err)
	}
	fmt.Println(page)
}

func intPtr(value int) *int {
	return &value
}
```

Destination catalog, status, sandbox, and jobs are zero-auth. Pass `sortedtravel.WithAPIKey(...)` or set `SORTED_TRAVEL_API_KEY` only when a host requires Bearer auth.

Interactive ranking, weather, and visa still go through MCP at `https://sorted.travel/mcp`.
