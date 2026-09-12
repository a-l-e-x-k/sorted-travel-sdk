"""Unit tests for the Sorted Travel REST client."""

import unittest

from sorted_travel.client import APIError, Client


class ClientTests(unittest.TestCase):
    """Cover URL building, auth headers, and error mapping."""

    def test_status_uses_public_path(self) -> None:
        """status() GETs /api/v1/status with no Authorization header."""
        captured = {}

        def transport(request, timeout):
            captured["request"] = request
            captured["timeout"] = timeout
            return 200, "application/json", '{"status":"OK"}'

        client = Client(transport=transport, env={})
        payload = client.status()
        self.assertEqual(payload, {"status": "OK"})
        self.assertEqual(
            captured["request"]["url"],
            "https://sorted.travel/api/v1/status",
        )
        self.assertEqual(captured["request"]["method"], "GET")
        self.assertNotIn("authorization", captured["request"]["headers"])

    def test_list_destinations_encodes_cursor(self) -> None:
        """list_destinations() passes cursor and limit as query params."""
        captured = {}

        def transport(request, timeout):
            del timeout
            captured["url"] = request["url"]
            return 200, "application/json", '{"data":[],"has_more":false}'

        client = Client(transport=transport, env={})
        client.list_destinations(cursor="abc+1", limit=5)
        self.assertIn("cursor=abc%2B1", captured["url"])
        self.assertIn("limit=5", captured["url"])

    def test_api_key_sets_bearer_header(self) -> None:
        """A configured key is sent as Authorization: Bearer."""
        captured = {}

        def transport(request, timeout):
            del timeout
            captured["headers"] = request["headers"]
            return 200, "application/json", "{}"

        client = Client(api_key="st_sandbox_test", transport=transport, env={})
        client.status()
        self.assertEqual(
            captured["headers"]["authorization"],
            "Bearer st_sandbox_test",
        )

    def test_non_2xx_raises_api_error(self) -> None:
        """HTTP errors surface as APIError with status and body."""

        def transport(request, timeout):
            del request, timeout
            return 429, "application/json", '{"error":"rate_limited"}'

        client = Client(transport=transport, env={})
        with self.assertRaises(APIError) as raised:
            client.status()
        self.assertEqual(raised.exception.status, 429)
        self.assertEqual(raised.exception.body, {"error": "rate_limited"})


if __name__ == "__main__":
    unittest.main()
