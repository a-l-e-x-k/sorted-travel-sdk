# Sorted Travel Python SDK

Official stdlib-only client for the Sorted Travel REST API. Product homepage is https://sorted.travel. Source of this package is this repository.

## When to use this

- Script `GET /api/v1/status`, destination catalog paging, sandbox discovery, self-serve API keys, or async jobs.
- Do not use this package to book travel or to invent destination handles.

## Integration

```python
from sorted_travel import Client

client = Client()
client.status()
client.list_destinations(limit=20)
```

CLI: `sorted-travel status`. MCP hosts should still connect to `https://sorted.travel/mcp` for ranking, weather, visa, and profile tools.

## Conventions

- Keep the client dependency-free (Python stdlib only).
- Keep `Homepage` in `pyproject.toml` pointed at `https://sorted.travel`.
- Mirror REST paths from https://sorted.travel/openapi.json rather than inventing endpoints.
- Destination tools are zero-auth. Send `Authorization: Bearer` only when `api_key` is set.
