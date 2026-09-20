# Sorted Travel SDKs

Official stdlib-only clients for the Sorted Travel REST API. Product homepage is https://sorted.travel. Source of this package is this repository.

## When to use this

- Script `GET /api/v1/status`, destination catalog paging, sandbox discovery, self-serve API keys, or async jobs.
- Do not use this package to book travel or to invent destination handles.

## Integration

Python:

```python
from sorted_travel import Client

client = Client()
client.status()
client.list_destinations(limit=20)
```

JavaScript:

```js
import { Client } from "sorted-travel"

const client = new Client()
await client.status()
await client.listDestinations({ limit: 20 })
```

Go:

```go
client := sortedtravel.NewClient()
client.Status(context.Background())
```

Ruby:

```ruby
require "sorted_travel"

client = SortedTravel::Client.new
client.status
```

CLI: `sorted-travel status`. MCP hosts should still connect to `https://sorted.travel/mcp` for ranking, weather, visa, and profile tools. This repo also ships an Agent Plugins v1 package (`plugin.json`, `skills/`, `mcp.json`) and a skills.sh entry (`npx skills add a-l-e-x-k/sorted-travel-sdk`).

## Conventions

- Keep the clients dependency-free (Python stdlib, JavaScript built-ins, Go stdlib).
- Keep `Homepage` in `pyproject.toml`, `homepage` in `js/package.json`, `homepage` / `homepage_uri` in `ruby/sorted_travel.gemspec`, and the Go README pointed at `https://sorted.travel`.
- Publish Go modules with tags like `go/v0.1.0` so `pkg.go.dev` indexes `github.com/a-l-e-x-k/sorted-travel-sdk/go/sortedtravel`.
- Mirror REST paths from https://sorted.travel/openapi.json rather than inventing endpoints.
- Destination tools are zero-auth. Send `Authorization: Bearer` only when `api_key` is set.
