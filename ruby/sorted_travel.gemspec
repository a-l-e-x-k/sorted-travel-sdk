# frozen_string_literal: true

require_relative "lib/sorted_travel/version"

Gem::Specification.new do |spec|
  spec.name = "sorted-travel"
  spec.version = SortedTravel::VERSION
  spec.authors = ["Sorted Travel"]
  spec.email = ["support@sorted.travel"]

  spec.summary = "Official Ruby SDK and CLI for the Sorted Travel REST API."
  spec.description = spec.summary
  spec.homepage = "https://sorted.travel"
  spec.license = "MIT"
  spec.licenses = ["MIT"]
  spec.required_ruby_version = ">= 3.1.0"

  spec.metadata = {
    "homepage_uri" => "https://sorted.travel",
    "source_code_uri" => "https://github.com/a-l-e-x-k/sorted-travel-sdk",
    "documentation_uri" => "https://sorted.travel/sdks.md",
    "rubygems_mfa_required" => "true",
  }

  spec.files = Dir.chdir(__dir__) do
    Dir["{lib,bin}/**/*", "README.md"]
  end
  spec.bindir = "bin"
  spec.executables = ["sorted-travel"]
  spec.require_paths = ["lib"]
end
