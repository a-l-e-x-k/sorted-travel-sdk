#!/usr/bin/env node

import { Client } from "./index.js"

function printJson(value) {
  process.stdout.write(`${JSON.stringify(value, null, 2)}\n`)
}

function parseArgs(argv) {
  const flags = {}
  const positionals = []
  for (let index = 0; index < argv.length; index += 1) {
    const token = argv[index]
    if (token === "--base-url" || token === "--api-key" || token === "--cursor" || token === "--limit") {
      flags[token.slice(2)] = argv[index + 1]
      index += 1
      continue
    }
    if (token.startsWith("-")) {
      throw new Error(`unknown flag ${token}`)
    }
    positionals.push(token)
  }
  return { flags, positionals }
}

async function main(argv = process.argv.slice(2)) {
  const { flags, positionals } = parseArgs(argv)
  const command = positionals[0]
  if (!command) {
    process.stderr.write(
      "usage: sorted-travel <status|sandbox|api-key|destinations|create-job|job> [args]\n",
    )
    return 2
  }
  const client = new Client({
    apiKey: flags["api-key"],
    baseUrl: flags["base-url"],
  })
  if (command === "status") {
    printJson(await client.status())
    return 0
  }
  if (command === "sandbox") {
    printJson(await client.sandbox())
    return 0
  }
  if (command === "api-key") {
    printJson(await client.createApiKey())
    return 0
  }
  if (command === "destinations") {
    const limit = flags.limit === undefined ? undefined : Number(flags.limit)
    printJson(await client.listDestinations({ cursor: flags.cursor, limit }))
    return 0
  }
  if (command === "create-job") {
    const operation = positionals[1]
    if (!operation) {
      throw new Error("create-job requires an operation")
    }
    const extra = {}
    if (flags.limit !== undefined) {
      extra.limit = Number(flags.limit)
    }
    printJson(await client.createJob(operation, extra))
    return 0
  }
  if (command === "job") {
    const jobId = positionals[1]
    if (!jobId) {
      throw new Error("job requires a job id")
    }
    printJson(await client.getJob(jobId))
    return 0
  }
  throw new Error(`unknown command ${command}`)
}

main().then(
  (code) => {
    process.exitCode = code
  },
  (error) => {
    process.stderr.write(`${error.message}\n`)
    process.exitCode = 1
  },
)
