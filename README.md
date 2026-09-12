# sorted-travel

Official Python SDK and CLI for the [Sorted Travel](https://sorted.travel) REST API.

This is a stdlib-only client for the public endpoints at `https://sorted.travel`. Destination ranking, facts, weather, and visa tools are a free tier with zero-auth access. An API key is optional.

- Homepage: [https://sorted.travel](https://sorted.travel)
- Developer portal: [https://sorted.travel/developers](https://sorted.travel/developers)
- SDK guide: [https://sorted.travel/sdks.md](https://sorted.travel/sdks.md)
- OpenAPI: [https://sorted.travel/openapi.json](https://sorted.travel/openapi.json)
- MCP: [https://sorted.travel/mcp](https://sorted.travel/mcp)
- PyPI: [https://pypi.org/project/sorted-travel/](https://pypi.org/project/sorted-travel/)

## Install

```sh
pip install sorted-travel
```

## Quickstart

```python
from sorted_travel import Client

client = Client()
print(client.status())
page = client.list_destinations(limit=20)
print(page["data"][0]["handle"], page["has_more"])
```

Destination tools need no key. If a host still requires one:

```python
client = Client(api_key="st_sandbox_...")  # or set SORTED_TRAVEL_API_KEY
```

Self-serve sandbox keys:

```python
payload = client.create_api_key()
```

## CLI

```sh
sorted-travel status
sorted-travel destinations --limit 5
sorted-travel sandbox
sorted-travel api-key
```

## Configuration

| Constructor arg | Environment variable | Default |
| --- | --- | --- |
| `api_key` | `SORTED_TRAVEL_API_KEY` | — |
| `base_url` | `SORTED_TRAVEL_BASE_URL` | `https://sorted.travel` |
| `timeout` | — | `30.0` seconds |

## Surface

- `status()` — `GET /api/v1/status`
- `list_destinations(cursor=, limit=)` — `GET /api/v1/destinations`
- `create_api_key()` — `POST /api/v1/api-keys`
- `sandbox()` — `GET /sandbox`
- `create_job(operation, **body)` — `POST /api/v1/jobs`
- `get_job(job_id)` — `GET /api/v1/jobs/{job_id}`
- `get(path, params=)` / `post(path, body=)` — raw REST escape hatch

Interactive ranking, weather, and visa still go through MCP at `https://sorted.travel/mcp`. This package does not wrap JSON-RPC.

## License

MIT. The Sorted Travel product remains under its own terms at [https://sorted.travel/terms](https://sorted.travel/terms).

## Publishing to PyPI

Trusted Publishing is already wired in `.github/workflows/publish.yml` (kept locally until the GitHub token has `workflow` scope). On PyPI, add a trusted publisher for project `sorted-travel`, repo `a-l-e-x-k/sorted-travel-sdk`, workflow `publish.yml`, environment `pypi`. Then tag a GitHub Release named `v0.1.0`.
