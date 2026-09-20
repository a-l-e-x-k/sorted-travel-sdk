# frozen_string_literal: true

require "json"
require "net/http"
require "uri"

module SortedTravel
  # Base error for this SDK.
  class SortedTravelError < StandardError
  end

  # Returned for non-2xx HTTP responses.
  class APIError < SortedTravelError
    attr_reader :status, :body

    # Store HTTP status and parsed response body.
    def initialize(status, body)
      @status = status
      @body = body
      super("HTTP #{status}: #{truncate_body(body)}")
    end

    private

    def truncate_body(value, limit = 300)
      text = value.is_a?(String) ? value : JSON.generate(value)
      return text if text.length <= limit

      "#{text[0, limit - 1]}…"
    end
  end

  # Thin client for the Sorted Travel REST API.
  class Client
    DEFAULT_BASE_URL = "https://sorted.travel"
    DEFAULT_TIMEOUT = 30.0
    USER_AGENT = "sorted-travel-ruby/#{SortedTravel::VERSION} (+https://sorted.travel)"

    # Build a client from options and optional environment variables.
    def initialize(api_key: nil, base_url: nil, timeout: DEFAULT_TIMEOUT, transport: nil, env: ENV)
      @api_key = api_key || env["SORTED_TRAVEL_API_KEY"]
      configured_base = base_url || env["SORTED_TRAVEL_BASE_URL"] || DEFAULT_BASE_URL
      @base_url = configured_base.sub(%r{/+\z}, "")
      @timeout = timeout
      @transport = transport || method(:default_transport)
    end

    # GET /api/v1/status.
    def status
      get("/api/v1/status")
    end

    # GET /sandbox.
    def sandbox
      get("/sandbox")
    end

    # GET /api/v1/destinations with cursor pagination.
    def list_destinations(cursor: nil, limit: nil)
      params = {}
      params["cursor"] = cursor unless cursor.nil?
      params["limit"] = limit unless limit.nil?
      get("/api/v1/destinations", params)
    end

    # POST /api/v1/api-keys.
    def create_api_key
      post("/api/v1/api-keys", {})
    end

    # POST /api/v1/jobs.
    def create_job(operation, **body)
      post("/api/v1/jobs", { operation: operation, **body })
    end

    # GET /api/v1/jobs/{job_id}.
    def get_job(job_id)
      encoded_job_id = URI.encode_www_form_component(job_id)
      get("/api/v1/jobs/#{encoded_job_id}")
    end

    # GET a host-relative REST path.
    def get(path, params = {})
      request("GET", path, params: params)
    end

    # POST JSON to a host-relative REST path.
    def post(path, body)
      request("POST", path, body: body)
    end

    private

    def request(method, path, params: {}, body: nil)
      raise ArgumentError, "REST paths must start with '/'" unless path.start_with?("/")

      url = @base_url + path
      unless params.empty?
        query = URI.encode_www_form(stringify_params(params))
        url = "#{url}?#{query}"
      end

      headers = {
        "User-Agent" => USER_AGENT,
        "Accept" => "application/json",
      }
      encoded_body = nil
      unless body.nil?
        headers["Content-Type"] = "application/json"
        encoded_body = JSON.generate(body)
      end
      headers["Authorization"] = "Bearer #{@api_key}" if @api_key

      status_code, _content_type, text = @transport.call(
        url: url,
        method: method,
        headers: headers,
        body: encoded_body,
        timeout: @timeout,
      )
      value = parse_body(text)
      raise APIError.new(status_code, value) if status_code < 200 || status_code >= 300

      value
    end

    def stringify_params(params)
      params.transform_values do |value|
        case value
        when true
          "true"
        when false
          "false"
        else
          value.to_s
        end
      end
    end

    def parse_body(text)
      return text if text.nil? || text.empty?

      JSON.parse(text)
    rescue JSON::ParserError
      text
    end

    def default_transport(url:, method:, headers:, body:, timeout:)
      uri = URI(url)
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = uri.scheme == "https"
      http.open_timeout = timeout
      http.read_timeout = timeout

      request_class = Net::HTTP.const_get(method.capitalize)
      http_request = request_class.new(uri)
      headers.each { |key, value| http_request[key] = value }
      http_request.body = body if body

      response = http.request(http_request)
      [response.code.to_i, response["Content-Type"], response.body]
    end
  end
end
