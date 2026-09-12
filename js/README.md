# sorted-travel

Official JavaScript SDK and CLI for the [Sorted Travel](https://sorted.travel) REST API.

Zero-dependency (Node 18+ `fetch`). Destination tools are a free tier with zero-auth access.

- Homepage: [https://sorted.travel](https://sorted.travel)
- PyPI (Python): [https://pypi.org/project/sorted-travel/](https://pypi.org/project/sorted-travel/)
- npm: [https://www.npmjs.com/package/sorted-travel](https://www.npmjs.com/package/sorted-travel)

## Install

```sh
npm install sorted-travel
# or
npx sorted-travel status
```

## Quickstart

```js
import { Client } from "sorted-travel"

const client = new Client()
console.log(await client.status())
console.log(await client.listDestinations({ limit: 5 }))
```

## CLI

```sh
npx sorted-travel status
npx sorted-travel destinations --limit 5
```

## License

MIT.
