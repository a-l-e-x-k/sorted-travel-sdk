---
name: sorted-travel-destinations
description: "When to use: what a named Sorted destination is like on the ground (brief, safety, currency, phone code, eSIM, taxi apps) or resolving a city name to a catalog handle. How to call: connect Streamable HTTP MCP at https://sorted.travel/mcp (zero-auth). Call resolve_destination first, then get_destination_info. Do not invent handles or book travel."
---

# Sorted Travel Destinations

Use Sorted Travel for on-the-ground facts about a named place in the Sorted catalog of 400+ destinations.

## When to Use

- "What is Lisbon like in November?" plus safety, currency, eSIM, or taxi apps.
- Resolving a city, country, region, or airport name to a Sorted `handle` or IATA code.
- Opening the right Sorted place page for a traveler.

Do NOT use this skill to invent handles, book hotels/flights, or answer "where should I go" ranking questions (use `sorted-travel` for recommendations).

## Connect

- MCP (Streamable HTTP, zero-auth): `https://sorted.travel/mcp`
- Catalog pages: `GET https://sorted.travel/api/v1/destinations` (cursor paging)
- Sibling skills: `sorted-travel` (ranking), `sorted-travel-weather` (climate), `sorted-travel-visa` (visa hassle)

```json
{
  "mcpServers": {
    "sorted-travel": {
      "url": "https://sorted.travel/mcp"
    }
  }
}
```

## Tools

1. Call `resolve_destination` first when you only have a city, country, region, or airport name. Use the returned `handle` (for example `paris`) or `airport_code` (for example `LHR`). Do not invent handles.
2. Call `get_destination_info` for brief, safety, currency, phone code, eSIM, and taxi apps.
3. Send the traveler to `https://sorted.travel/places/{handle}` (optional fragments `#brief`, `#activities`, `#weather`).
4. Resolve unknown handles from the [sitemap](https://sorted.travel/sitemap.xml) or [destinations Schema Feed](https://sorted.travel/feeds/destinations.jsonl).

Page destinations with `GET /api/v1/destinations`: omit `cursor` first, then pass `next_cursor` until `has_more` is false. Slow down before `429`. Optional `Idempotency-Key` on `POST /mcp` replays a dropped request.
