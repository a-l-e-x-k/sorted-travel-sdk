"""Command-line entry point for the Sorted Travel REST SDK."""

import argparse
import json
import sys

from sorted_travel.client import Client


def _print_json(value: object) -> int:
    """Write a JSON payload to stdout."""
    json.dump(value, sys.stdout, indent=2)
    sys.stdout.write("\n")
    return 0


def main(argv: list[str] | None = None) -> int:
    """Run the sorted-travel CLI."""
    parser = argparse.ArgumentParser(
        prog="sorted-travel",
        description="Official CLI for the Sorted Travel REST API.",
    )
    parser.add_argument(
        "--base-url",
        default=None,
        help="API origin. Defaults to https://sorted.travel or SORTED_TRAVEL_BASE_URL.",
    )
    parser.add_argument(
        "--api-key",
        default=None,
        help="Optional Bearer token. Defaults to SORTED_TRAVEL_API_KEY.",
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    subparsers.add_parser("status", help="GET /api/v1/status")
    subparsers.add_parser("sandbox", help="GET /sandbox")
    subparsers.add_parser("api-key", help="POST /api/v1/api-keys")

    destinations = subparsers.add_parser(
        "destinations",
        help="GET /api/v1/destinations",
    )
    destinations.add_argument("--cursor", default=None)
    destinations.add_argument("--limit", type=int, default=None)

    job_create = subparsers.add_parser("create-job", help="POST /api/v1/jobs")
    job_create.add_argument("operation")
    job_create.add_argument("--limit", type=int, default=None)

    job_get = subparsers.add_parser("job", help="GET /api/v1/jobs/{job_id}")
    job_get.add_argument("job_id")

    args = parser.parse_args(argv)
    client = Client(api_key=args.api_key, base_url=args.base_url)

    if args.command == "status":
        return _print_json(client.status())
    if args.command == "sandbox":
        return _print_json(client.sandbox())
    if args.command == "api-key":
        return _print_json(client.create_api_key())
    if args.command == "destinations":
        return _print_json(
            client.list_destinations(cursor=args.cursor, limit=args.limit)
        )
    if args.command == "create-job":
        extra = {}
        if args.limit is not None:
            extra["limit"] = args.limit
        return _print_json(client.create_job(args.operation, **extra))
    if args.command == "job":
        return _print_json(client.get_job(args.job_id))
    parser.error(f"unknown command {args.command}")
    return 2


if __name__ == "__main__":
    raise SystemExit(main())
