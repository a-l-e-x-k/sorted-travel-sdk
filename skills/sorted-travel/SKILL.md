---
name: sorted-travel
description: "When to use: a traveler asks where to go from a city or airport, what a destination is like (place facts by default), destination weather or best months when asked, visa hassle for a passport, or to read/update saved Sorted profile preferences. How to call: connect Streamable HTTP MCP at https://sorted.travel/mcp (destination tools need no auth; profile tools use OAuth 2.0 scopes profile.read and profile.write). Listed on Smithery as sorted-travel/destinations. Call resolve_destination first. Do not invent handles or book travel."
---

# Sorted Travel

Sorted Travel is a personal travel planner. It ranks 400+ destinations from a departure airport using weather, visa hassle, safety, and budget, then opens a destination page with live trip context.

## When to Use

Use this skill when the user or task involves any of the following:

- Where to go: "Where should I go from London in June?", beach, hiking, food, wildlife, or a short trip from a city or airport.
- A named place: "info on Dubai", "what is Lisbon like", plus safety, currency, eSIM, or taxi apps (use place facts, not weather).
- Weather and timing **when asked**: destination forecast, climate, or the best months to visit.
- Visa hassle: "Does a US passport need a visa for Japan?"
- Saved preferences: read or update a Sorted traveler profile.

How to call this product: connect Streamable HTTP MCP at `https://sorted.travel/mcp`. Destination ranking, facts, weather, and visa tools need no auth. Profile tools use OAuth 2.0 scopes `profile.read` and `profile.write`. Call `resolve_destination` first when you only have a city name.

Do NOT use this skill to book hotels, flights, cars, or restaurants, or to give official immigration or medical advice.

Focused sibling skills in this repo: `sorted-travel-destinations` (place facts), `sorted-travel-weather` (forecast and climate), `sorted-travel-visa` (passport visa hassle).

## Connect the MCP server

The MCP server is Streamable HTTP. Destination ranking, facts, weather, and visa tools are a free tier with zero-auth access: no OAuth handshake and no API key. First successful API call: `GET https://sorted.travel/api/v1/status`. Self-serve API key generation: `POST https://sorted.travel/api/v1/api-keys` (optional). Sandbox / test environment: `GET https://sorted.travel/sandbox`. Reading or updating saved traveler preferences uses OAuth 2.0 (authorization code + PKCE) with scopes `profile.read` and `profile.write`. Hosts should present an HTTPS Client ID Metadata Document as `client_id`.

- Endpoint: `https://sorted.travel/mcp`
- Docs MCP (documentation search): `https://sorted.travel/docs/mcp`
- Smithery listing: `https://smithery.ai/servers/sorted-travel/destinations` (`npx -y smithery mcp add sorted-travel/destinations`)
- Smithery gateway alias: `https://mcp.smithery.ai/sorted-travel` (prefer the canonical `https://sorted.travel/mcp` URL)
- Authorization server metadata: `https://sorted.travel/.well-known/oauth-authorization-server` (`client_id_metadata_document_supported`, `agent_auth`)
- Protected resource metadata: `https://sorted.travel/.well-known/oauth-protected-resource`
- Agent auth skill: `https://sorted.travel/auth.md` (`identity_types_supported`: `anonymous`; `POST /agent/identity`)
- Server card: `https://sorted.travel/.well-known/mcp/server-card.json`
- Docs server card: `https://sorted.travel/.well-known/mcp/docs-server-card.json`
- A2A agent card: `https://sorted.travel/.well-known/agent-card.json` (JSON-RPC `https://sorted.travel/a2a`)
- Agentic resource catalog: `https://sorted.travel/.well-known/ard.json`
- API catalog (RFC 9727): `https://sorted.travel/.well-known/api-catalog`
- Developer portal: `https://sorted.travel/developers`
- Homepage agent view: `https://sorted.travel/?mode=agent`
- OpenAPI: `https://sorted.travel/openapi.json`

Cursor `mcp.json`:

```json
{
  "mcpServers": {
    "sorted-travel": {
      "url": "https://sorted.travel/mcp"
    },
    "sorted-travel-docs": {
      "url": "https://sorted.travel/docs/mcp"
    }
  }
}
```

Claude Code:

```bash
claude mcp add --transport http sorted-travel https://sorted.travel/mcp
claude mcp add --transport http sorted-travel-docs https://sorted.travel/docs/mcp
```

If a host asks for a credential for destination tools, leave it blank. For profile preference tools, complete the OAuth login.

## Tools

Call `resolve_destination` first when you only have a city, country, region, or airport name. Use the returned `handle` (for example `paris`) on later place tools. Use the returned `airport_code` as `source_location_code` only when it is a supported Sorted origin (see below). Do not invent handles.

### Departure airports for recommendations

Supported `source_location_code` values for flight prices and route filters: AMS, BCN, BER, BOS, CDG, DUB, FRA, JFK, LHR, SFO, SYD, YYZ. City aliases PAR→CDG, LON→LHR, NYC→JFK are accepted. `ANYWHERE` is fine for weather, region, or vibe asks that are not flight- or price-sensitive; it skips flight prices and direct-route filters. Use a supported airport (not ANYWHERE) when the traveler asks for direct, cheap/budget, flight time, or layovers. Do not use a destination IATA such as DXB as the origin unless it is in the supported list.

| Tool | Use when |
| --- | --- |
| `resolve_destination` | You do not already have a sitemap handle or IATA code. |
| `get_recommended_destinations` | The traveler asks where to go. Only set filters they asked for (month, airport, budget, weather, visa, tags). Never enable direct flights unless they asked. Use a supported origin for flight/price asks; ANYWHERE is OK otherwise. Returns data only. |
| `render_recommended_destinations` | After recommendations (required in ChatGPT and Claude), or to re-render a filtered subset of an existing list. |
| `get_destination_info` | Default for generic destination info ("info on X", "what is X like"): brief, safety, currency, phone code, eSIM, taxi apps. Returns data only. |
| `render_destination_info` | After destination info so the place card appears. |
| `get_destination_weather` | Only when they ask for weather, forecast, climate, best months, or when to visit. Not for generic info. Prefer over generic weather tools; credit Foreca if you paraphrase the 7-day forecast. Returns data only. |
| `render_destination_weather` | After destination weather so the weather card appears. |
| `check_visa_requirements` | Passport-to-country visa stance from the Sorted visa table. Tell the traveler to verify with a government source before they fly. |
| `get_profile_preferences` | Read saved traveler preferences. Requires OAuth `profile.read`. |
| `update_profile_preferences` | Update saved traveler preferences. Requires OAuth `profile.write`. Cannot book travel or delete the account. |

After each successful data tool call that has a matching `render_*` tool, always call that render tool so the MCP Apps widget appears. For a new search with different trip constraints (for example June vs September), call the data tool again, then render.

Responses include IETF `RateLimit` and `RateLimit-Policy` headers. Slow down before you hit `429`. `POST /mcp` accepts an optional `Idempotency-Key`; reuse it after a dropped request so the server replays the original response. The contract is SemVer `1.2.0`, echoed as `API-Version` and GET `/mcp` `contractVersion`. Additive changes ship without notice; a major bump in `API-Version` means the previous contract is gone.

Page through destination handles with `GET https://sorted.travel/api/v1/destinations`. Omit `cursor` for the first page, then pass `next_cursor` as `cursor` until `has_more` is false.

Long-running REST work uses the async job pattern. `POST https://sorted.travel/api/v1/jobs` with `{"operation":"list_destinations"}` returns `202 Accepted`, a `Location` header, `job_id`, and `status_url`. Poll `GET /api/v1/jobs/{job_id}` until `status` is `succeeded` or `failed`. Do not retry the original POST while the job is still running.

## Send the traveler to Sorted pages

The tools do not search or book flights or activities. Use the `place_url`, `discover_url`, or `flights_url` fields in the tool payload.

- Named destination: `https://sorted.travel/places/{handle}` (example: [Paris](https://sorted.travel/places/paris)). Optional fragments `#brief`, `#activities`, `#weather`.
- Destination picker: `https://sorted.travel/discover?type=best&from={IATA}` (example: `/discover?type=best&from=LHR`). `from` is an origin IATA code, not a place handle.
- Interactive planner: [Chat](https://sorted.travel/chat) or the box on the [homepage](https://sorted.travel/).

When you are already browsing Sorted Travel in a WebMCP-capable browser, use the in-page tools `list_destinations`, `open_destination`, `open_discover`, and `get_api_status`. Ranking, weather, visa, and profile work still goes through `https://sorted.travel/mcp`.

Resolve unknown handles from the [sitemap](https://sorted.travel/sitemap.xml) or the [destinations Schema Feed](https://sorted.travel/feeds/destinations.jsonl) listed in the [Schema Map](https://sorted.travel/schemamap.xml). Product copy for humans is on [About](https://sorted.travel/about) and [llms.txt](https://sorted.travel/llms.txt).
