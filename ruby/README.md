# sorted-travel Ruby SDK

Official Ruby client and CLI for the [Sorted Travel](https://sorted.travel) REST API.

- Homepage: [https://sorted.travel](https://sorted.travel)
- Developer portal: [https://sorted.travel/developers](https://sorted.travel/developers)
- SDK guide: [https://sorted.travel/sdks.md](https://sorted.travel/sdks.md)
- Repository: [https://github.com/a-l-e-x-k/sorted-travel-sdk](https://github.com/a-l-e-x-k/sorted-travel-sdk)

## Install

```sh
gem install sorted-travel
```

## Quickstart

```ruby
require "sorted_travel"

client = SortedTravel::Client.new
client.status
client.list_destinations(limit: 20)
```

Destination catalog, status, sandbox, and jobs are zero-auth. Pass `api_key:` or set `SORTED_TRAVEL_API_KEY` only when a host requires Bearer auth.

Interactive ranking, weather, and visa still go through MCP at `https://sorted.travel/mcp`.
