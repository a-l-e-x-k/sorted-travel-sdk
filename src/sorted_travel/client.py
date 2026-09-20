"""Stdlib HTTP client for Sorted Travel public REST endpoints."""

import json
import os
import urllib.error
import urllib.parse
import urllib.request

__version__ = "0.1.1"

DEFAULT_BASE_URL = "https://sorted.travel"
DEFAULT_TIMEOUT = 30.0
USER_AGENT = f"sorted-travel-python/{__version__} (+https://sorted.travel)"


class SortedTravelError(Exception):
    """Base class for errors raised by this SDK."""


class APIError(SortedTravelError):
    """A REST or transport-level failure (non-2xx HTTP response)."""

    def __init__(self, status: int, body: object) -> None:
        """Store the HTTP status and parsed body."""
        self.status = status
        self.body = body
        super().__init__(f"HTTP {status}: {_truncate(body)}")


def _truncate(value: object, limit: int = 300) -> str:
    """Return a short string for exception messages."""
    text = value if isinstance(value, str) else json.dumps(value)
    if len(text) <= limit:
        return text
    return text[: limit - 1] + "…"


def _stringify(value: object) -> str:
    """Encode a query value for urllib."""
    if value is True:
        return "true"
    if value is False:
        return "false"
    return str(value)


def _default_transport(request: dict, timeout: float) -> tuple[int, str, str]:
    """Perform an HTTP request with urllib.

    ``request`` has keys url, method, headers, and optional body bytes.
    Non-2xx responses are returned, not raised.
    """
    req = urllib.request.Request(
        request["url"],
        data=request.get("body"),
        headers=request.get("headers") or {},
        method=request.get("method", "GET"),
    )
    try:
        with urllib.request.urlopen(req, timeout=timeout) as response:
            return (
                response.status,
                response.headers.get("content-type", ""),
                response.read().decode("utf-8", "replace"),
            )
    except urllib.error.HTTPError as err:
        return (
            err.code,
            err.headers.get("content-type", "") if err.headers else "",
            err.read().decode("utf-8", "replace"),
        )


def parse_body(text: str, content_type: str = "") -> object:
    """Decode a JSON body, or return the raw text when it is not JSON."""
    del content_type
    if not text:
        return text
    try:
        return json.loads(text)
    except ValueError:
        return text


class Client:
    """Thin client for the Sorted Travel REST API.

    Destination catalog, status, sandbox, and jobs are zero-auth. Pass
    ``api_key`` only when a host requires ``Authorization: Bearer``.
    ``transport`` is injectable for tests: ``(request_dict, timeout) ->
    (status, content_type, text)``.
    """

    def __init__(
        self,
        api_key: str | None = None,
        base_url: str | None = None,
        timeout: float = DEFAULT_TIMEOUT,
        transport=None,
        env: dict[str, str] | None = None,
    ) -> None:
        """Load credentials and base URL from arguments or the environment."""
        environ = os.environ if env is None else env
        self.api_key = api_key or environ.get("SORTED_TRAVEL_API_KEY")
        configured_base = base_url or environ.get("SORTED_TRAVEL_BASE_URL") or DEFAULT_BASE_URL
        self.base_url = configured_base.rstrip("/")
        self.timeout = timeout
        self._transport = transport or _default_transport

    def status(self) -> object:
        """GET /api/v1/status."""
        return self.get("/api/v1/status")

    def sandbox(self) -> object:
        """GET /sandbox."""
        return self.get("/sandbox")

    def list_destinations(
        self,
        cursor: str | None = None,
        limit: int | None = None,
    ) -> object:
        """GET /api/v1/destinations with cursor pagination."""
        params: dict[str, object] = {}
        if cursor is not None:
            params["cursor"] = cursor
        if limit is not None:
            params["limit"] = limit
        return self.get("/api/v1/destinations", params=params)

    def create_api_key(self) -> object:
        """POST /api/v1/api-keys and return a sandbox credential."""
        return self.post("/api/v1/api-keys", body={})

    def create_job(self, operation: str, **body: object) -> object:
        """POST /api/v1/jobs and return the 202 Accepted payload."""
        payload = {"operation": operation, **body}
        return self.post("/api/v1/jobs", body=payload)

    def get_job(self, job_id: str) -> object:
        """GET /api/v1/jobs/{job_id}."""
        encoded_job_id = urllib.parse.quote(job_id, safe="")
        return self.get(f"/api/v1/jobs/{encoded_job_id}")

    def get(self, path: str, params: dict[str, object] | None = None) -> object:
        """GET a host-relative REST path."""
        return self._request("GET", path, params=params)

    def post(self, path: str, body: object | None = None) -> object:
        """POST JSON to a host-relative REST path."""
        return self._request("POST", path, body=body)

    def _headers(self, extra: dict[str, str] | None = None) -> dict[str, str]:
        """Build request headers, including an optional Bearer token."""
        headers = {
            "user-agent": USER_AGENT,
            "accept": "application/json",
        }
        if extra:
            headers.update(extra)
        if self.api_key:
            headers["authorization"] = f"Bearer {self.api_key}"
        return headers

    def _request(
        self,
        method: str,
        path: str,
        params: dict[str, object] | None = None,
        body: object | None = None,
    ) -> object:
        """Send one REST request and raise APIError on non-2xx."""
        if not path.startswith("/"):
            raise ValueError("REST paths must start with '/'")
        url = self.base_url + path
        if params:
            encoded = urllib.parse.urlencode(
                {key: _stringify(value) for key, value in params.items()}
            )
            url += "?" + encoded
        headers = self._headers()
        encoded_body = None
        if body is not None:
            headers["content-type"] = "application/json"
            encoded_body = json.dumps(body).encode("utf-8")
        status, content_type, text = self._transport(
            {
                "url": url,
                "method": method,
                "headers": headers,
                "body": encoded_body,
            },
            self.timeout,
        )
        value = parse_body(text, content_type)
        if status < 200 or status >= 300:
            raise APIError(status, value)
        return value
