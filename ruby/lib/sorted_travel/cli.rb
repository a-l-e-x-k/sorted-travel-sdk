# frozen_string_literal: true

require "json"
require "optparse"

module SortedTravel
  # Command-line entry point for the Sorted Travel REST SDK.
  module CLI
    module_function

    # Run the sorted-travel CLI and return an exit code.
    def run(argv = ARGV)
      options = { api_key: nil, base_url: nil }
      parser = OptionParser.new do |opts|
        opts.banner = "Usage: sorted-travel [options] COMMAND"
        opts.on("--base-url URL", "API origin (default https://sorted.travel or SORTED_TRAVEL_BASE_URL)") do |value|
          options[:base_url] = value
        end
        opts.on("--api-key KEY", "Optional Bearer token (default SORTED_TRAVEL_API_KEY)") do |value|
          options[:api_key] = value
        end
      end

      positional = parser.order(argv)
      parser.parse!(argv)
      command = positional.shift
      if command.nil?
        warn parser
        return 2
      end

      client = Client.new(api_key: options[:api_key], base_url: options[:base_url])

      case command
      when "status"
        print_json(client.status)
      when "sandbox"
        print_json(client.sandbox)
      when "api-key"
        print_json(client.create_api_key)
      when "destinations"
        destination_options = { cursor: nil, limit: nil }
        OptionParser.new do |opts|
          opts.on("--cursor VALUE") { |value| destination_options[:cursor] = value }
          opts.on("--limit INTEGER", Integer) { |value| destination_options[:limit] = value }
        end.parse!(positional)
        print_json(client.list_destinations(**destination_options))
      when "create-job"
        operation = positional.shift
        if operation.nil?
          warn "missing operation"
          return 2
        end
        job_options = {}
        OptionParser.new do |opts|
          opts.on("--limit INTEGER", Integer) { |value| job_options[:limit] = value }
        end.parse!(positional)
        print_json(client.create_job(operation, **job_options))
      when "job"
        job_id = positional.shift
        if job_id.nil?
          warn "missing job_id"
          return 2
        end
        print_json(client.get_job(job_id))
      else
        warn "unknown command #{command}"
        return 2
      end
      0
    rescue OptionParser::ParseError => err
      warn err.message
      2
    end

    def print_json(value)
      puts JSON.pretty_generate(value)
    end
  end
end
