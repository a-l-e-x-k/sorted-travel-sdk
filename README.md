# sorted-travel

[![skills.sh](https://skills.sh/b/a-l-e-x-k/sorted-travel-sdk)](https://skills.sh/a-l-e-x-k/sorted-travel-sdk)

Official Python and JavaScript SDK and CLI for the [Sorted Travel](https://sorted.travel) REST API.

This repository publishes:

- Python: [`sorted-travel` on PyPI](https://pypi.org/project/sorted-travel/) (`pip install sorted-travel`)
- JavaScript: [`sorted-travel` on npm](https://www.npmjs.com/package/sorted-travel/) (`npm install sorted-travel`, source in [`js/`](js/))
- Agent skill: [`skills/sorted-travel/SKILL.md`](skills/sorted-travel/SKILL.md) on [skills.sh](https://skills.sh/a-l-e-x-k/sorted-travel-sdk/sorted-travel)
- Agent Plugin: [`plugin.json`](plugin.json), [`skills/`](skills/), and [`mcp.json`](mcp.json) (Agent Plugins v1)

This is a stdlib-only Python client and a zero-dependency JavaScript client for the public endpoints at `https://sorted.travel`. Destination ranking, facts, weather, and visa tools are a free tier with zero-auth access. An API key is optional.

- Homepage: [https://sorted.travel](https://sorted.travel)
- Developer portal: [https://sorted.travel/developers](https://sorted.travel/developers)
- SDK guide: [https://sorted.travel/sdks.md](https://sorted.travel/sdks.md)
- OpenAPI: [https://sorted.travel/openapi.json](https://sorted.travel/openapi.json)
- MCP: [https://sorted.travel/mcp](https://sorted.travel/mcp)
- Smithery: [sorted-travel/destinations](https://smithery.ai/servers/sorted-travel/destinations)
- PyPI: [https://pypi.org/project/sorted-travel/](https://pypi.org/project/sorted-travel/)
- npm: [https://www.npmjs.com/package/sorted-travel/](https://www.npmjs.com/package/sorted-travel/)

## Agent skill (skills.sh)

Install the Sorted Travel agent skill for Cursor, Claude Code, Codex, and other supported agents:

```sh
npx skills add a-l-e-x-k/sorted-travel-sdk
```

Non-interactive install for Cursor:

```sh
npx skills add a-l-e-x-k/sorted-travel-sdk -y -a cursor
```

The skill teaches agents when to use Sorted Travel, how to connect MCP at `https://sorted.travel/mcp`, and which tools to call. Source: [`skills/sorted-travel/SKILL.md`](skills/sorted-travel/SKILL.md).

Connect MCP directly with [`mcp.json`](mcp.json) or Smithery:

```sh
npx -y smithery mcp add sorted-travel/destinations
```

## Install SDKs

Python:

```sh
pip install sorted-travel
```

JavaScript:

```sh
npm install sorted-travel
```

## Quickstart (Python)

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

## Quickstart (JavaScript)

```js
import { Client } from "sorted-travel"

const client = new Client()
await client.status()
await client.listDestinations({ limit: 20 })
```

## CLI

```sh
sorted-travel status
sorted-travel destinations --limit 5
sorted-travel sandbox
sorted-travel api-key
```

JavaScript CLI via npm:

```sh
npx sorted-travel status
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

Interactive ranking, weather, and visa still go through MCP at `https://sorted.travel/mcp`. These SDK packages do not wrap JSON-RPC.

## License

MIT. The Sorted Travel product remains under its own terms at [https://sorted.travel/terms](https://sorted.travel/terms).
