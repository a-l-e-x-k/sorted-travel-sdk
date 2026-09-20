---
name: sorted-travel-visa
description: "When to use: passport-to-country visa hassle for a Sorted destination (for example does a US passport need a visa for Japan). How to call: connect Streamable HTTP MCP at https://sorted.travel/mcp (zero-auth). Call resolve_destination first when you only have a city or country name, then check_visa_requirements. Tell the traveler to verify with a government source before they fly. Do not give official immigration advice."
---

# Sorted Travel Visa

Use Sorted Travel to look up passport-to-country visa stance from the Sorted visa table for destinations in the catalog.

## When to Use

- "Does a US passport need a visa for Japan?"
- Visa hassle when ranking or comparing Sorted destinations.
- Quick orientation before the traveler checks a government source.

Do NOT use this skill as official immigration, entry, or medical advice. Always tell the traveler to verify with a government source before they fly. Do not book travel.

## Connect

- MCP (Streamable HTTP, zero-auth): `https://sorted.travel/mcp`
- Sibling skills: `sorted-travel` (ranking), `sorted-travel-destinations` (place facts), `sorted-travel-weather` (climate)

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

1. Call `resolve_destination` when you only have a city or country name. Use the returned country / handle context. Do not invent handles.
2. Call `check_visa_requirements` with the traveler's passport country and the destination country from Sorted.
3. State clearly that the result is Sorted's visa-hassle table, not a government ruling.
4. Link the place page when useful: `https://sorted.travel/places/{handle}`.

Slow down before `429` (`RateLimit` / `RateLimit-Policy` headers). Optional `Idempotency-Key` on `POST /mcp` replays a dropped request.
