# frozen_string_literal: true

require "minitest/autorun"

require_relative "../lib/sorted_travel"

class ClientTest < Minitest::Test
  def test_status_uses_public_path
    captured = nil
    client = SortedTravel::Client.new(
      env: {},
      transport: lambda do |url:, method:, headers:, body:, timeout:|
        captured = { url: url, method: method, headers: headers, body: body }
        [200, "application/json", '{"status":"OK"}']
      end,
    )
    payload = client.status
    assert_equal({ "status" => "OK" }, payload)
    assert_equal("https://sorted.travel/api/v1/status", captured[:url])
    assert_equal("GET", captured[:method])
    assert_nil(captured[:headers]["Authorization"])
  end

  def test_list_destinations_encodes_cursor
    captured_url = nil
    client = SortedTravel::Client.new(
      env: {},
      transport: lambda do |url:, method:, headers:, body:, timeout:|
        captured_url = url
        [200, "application/json", '{"data":[],"has_more":false}']
      end,
    )
    client.list_destinations(cursor: "abc+1", limit: 5)
    assert_match(/cursor=abc%2B1/, captured_url)
    assert_match(/limit=5/, captured_url)
  end

  def test_api_key_sets_bearer_header
    captured_headers = nil
    client = SortedTravel::Client.new(
      api_key: "st_sandbox_test",
      env: {},
      transport: lambda do |url:, method:, headers:, body:, timeout:|
        captured_headers = headers
        [200, "application/json", "{}"]
      end,
    )
    client.status
    assert_equal("Bearer st_sandbox_test", captured_headers["Authorization"])
  end

  def test_non2xx_raises_api_error
    client = SortedTravel::Client.new(
      env: {},
      transport: lambda do |url:, method:, headers:, body:, timeout:|
        [429, "application/json", '{"error":"rate_limited"}']
      end,
    )
    error = assert_raises(SortedTravel::APIError) { client.status }
    assert_equal(429, error.status)
    assert_equal({ "error" => "rate_limited" }, error.body)
  end
end
