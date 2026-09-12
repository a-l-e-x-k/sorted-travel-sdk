import assert from "node:assert/strict"
import test from "node:test"

import { APIError, Client } from "./index.js"

test("status GETs /api/v1/status with no Authorization header", async () => {
  let captured
  const client = new Client({
    env: {},
    transport: async (request) => {
      captured = request
      return { status: 200, contentType: "application/json", text: '{"status":"OK"}' }
    },
  })
  const payload = await client.status()
  assert.deepEqual(payload, { status: "OK" })
  assert.equal(captured.url, "https://sorted.travel/api/v1/status")
  assert.equal(captured.method, "GET")
  assert.equal(captured.headers.authorization, undefined)
})

test("listDestinations encodes cursor query params", async () => {
  let capturedUrl
  const client = new Client({
    env: {},
    transport: async (request) => {
      capturedUrl = request.url
      return { status: 200, contentType: "application/json", text: '{"data":[]}' }
    },
  })
  await client.listDestinations({ cursor: "abc+1", limit: 5 })
  assert.match(capturedUrl, /cursor=abc%2B1/)
  assert.match(capturedUrl, /limit=5/)
})

test("apiKey sets Bearer header", async () => {
  let capturedHeaders
  const client = new Client({
    apiKey: "st_sandbox_test",
    env: {},
    transport: async (request) => {
      capturedHeaders = request.headers
      return { status: 200, contentType: "application/json", text: "{}" }
    },
  })
  await client.status()
  assert.equal(capturedHeaders.authorization, "Bearer st_sandbox_test")
})

test("non-2xx raises APIError", async () => {
  const client = new Client({
    env: {},
    transport: async () => ({
      status: 429,
      contentType: "application/json",
      text: '{"error":"rate_limited"}',
    }),
  })
  await assert.rejects(
    () => client.status(),
    (error) => {
      assert.equal(error instanceof APIError, true)
      assert.equal(error.status, 429)
      assert.deepEqual(error.body, { error: "rate_limited" })
      return true
    },
  )
})
