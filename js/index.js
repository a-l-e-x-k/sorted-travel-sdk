const DEFAULT_BASE_URL = "https://sorted.travel"
const DEFAULT_TIMEOUT_MS = 30000
const VERSION = "0.1.0"
const USER_AGENT = `sorted-travel-js/${VERSION} (+https://sorted.travel)`

export class SortedTravelError extends Error {}

export class APIError extends SortedTravelError {
  constructor(status, body) {
    super(`HTTP ${status}`)
    this.status = status
    this.body = body
  }
}

function stringifyQueryValue(value) {
  if (value === true) {
    return "true"
  }
  if (value === false) {
    return "false"
  }
  return String(value)
}

async function defaultTransport(request, timeoutMs) {
  const controller = new AbortController()
  const timer = setTimeout(() => controller.abort(), timeoutMs)
  try {
    const response = await fetch(request.url, {
      method: request.method,
      headers: request.headers,
      body: request.body,
      signal: controller.signal,
    })
    const text = await response.text()
    return {
      status: response.status,
      contentType: response.headers.get("content-type") || "",
      text,
    }
  } finally {
    clearTimeout(timer)
  }
}

function parseBody(text) {
  if (!text) {
    return text
  }
  try {
    return JSON.parse(text)
  } catch {
    return text
  }
}

export class Client {
  constructor({
    apiKey,
    baseUrl,
    timeoutMs = DEFAULT_TIMEOUT_MS,
    transport,
    env = process.env,
  } = {}) {
    this.apiKey = apiKey || env.SORTED_TRAVEL_API_KEY
    const configuredBase = baseUrl || env.SORTED_TRAVEL_BASE_URL || DEFAULT_BASE_URL
    this.baseUrl = configuredBase.replace(/\/+$/, "")
    this.timeoutMs = timeoutMs
    this.transport = transport || defaultTransport
  }

  status() {
    return this.get("/api/v1/status")
  }

  sandbox() {
    return this.get("/sandbox")
  }

  listDestinations({ cursor, limit } = {}) {
    const params = {}
    if (cursor !== undefined) {
      params.cursor = cursor
    }
    if (limit !== undefined) {
      params.limit = limit
    }
    return this.get("/api/v1/destinations", params)
  }

  createApiKey() {
    return this.post("/api/v1/api-keys", {})
  }

  createJob(operation, body = {}) {
    return this.post("/api/v1/jobs", { operation, ...body })
  }

  getJob(jobId) {
    return this.get(`/api/v1/jobs/${encodeURIComponent(jobId)}`)
  }

  get(path, params = {}) {
    return this.request("GET", path, { params })
  }

  post(path, body) {
    return this.request("POST", path, { body })
  }

  async request(method, path, { params, body } = {}) {
    if (!path.startsWith("/")) {
      throw new Error("REST paths must start with '/'")
    }
    const url = new URL(this.baseUrl + path)
    for (const [key, value] of Object.entries(params || {})) {
      url.searchParams.set(key, stringifyQueryValue(value))
    }
    const headers = {
      "user-agent": USER_AGENT,
      accept: "application/json",
    }
    let encodedBody
    if (body !== undefined) {
      headers["content-type"] = "application/json"
      encodedBody = JSON.stringify(body)
    }
    if (this.apiKey) {
      headers.authorization = `Bearer ${this.apiKey}`
    }
    const result = await this.transport(
      {
        url: url.toString(),
        method,
        headers,
        body: encodedBody,
      },
      this.timeoutMs,
    )
    const value = parseBody(result.text)
    if (result.status < 200 || result.status >= 300) {
      throw new APIError(result.status, value)
    }
    return value
  }
}

export { DEFAULT_BASE_URL, USER_AGENT, VERSION }
