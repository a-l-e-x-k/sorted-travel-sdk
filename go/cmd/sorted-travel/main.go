package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/a-l-e-x-k/sorted-travel-sdk/go/sortedtravel"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(argv []string) int {
	flags, positionals, err := parseArgs(argv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	if len(positionals) == 0 {
		fmt.Fprintln(os.Stderr, "usage: sorted-travel <status|sandbox|api-key|destinations|create-job|job> [args]")
		return 2
	}

	options := []sortedtravel.ClientOption{}
	if flags["base-url"] != "" {
		options = append(options, sortedtravel.WithBaseURL(flags["base-url"]))
	}
	if flags["api-key"] != "" {
		options = append(options, sortedtravel.WithAPIKey(flags["api-key"]))
	}
	client := sortedtravel.NewClient(options...)
	command := positionals[0]

	switch command {
	case "status":
		return printJSON(client.Status(context.Background()))
	case "sandbox":
		return printJSON(client.Sandbox(context.Background()))
	case "api-key":
		return printJSON(client.CreateAPIKey(context.Background()))
	case "destinations":
		var cursor *string
		if cursorValue := flags["cursor"]; cursorValue != "" {
			cursor = &cursorValue
		}
		var limit *int
		if flags["limit"] != "" {
			parsedLimit, parseErr := strconv.Atoi(flags["limit"])
			if parseErr != nil {
				fmt.Fprintln(os.Stderr, parseErr.Error())
				return 1
			}
			limit = &parsedLimit
		}
		return printJSON(client.ListDestinations(context.Background(), cursor, limit))
	case "create-job":
		if len(positionals) < 2 {
			fmt.Fprintln(os.Stderr, "create-job requires an operation")
			return 1
		}
		body := map[string]interface{}{}
		if flags["limit"] != "" {
			parsedLimit, parseErr := strconv.Atoi(flags["limit"])
			if parseErr != nil {
				fmt.Fprintln(os.Stderr, parseErr.Error())
				return 1
			}
			body["limit"] = parsedLimit
		}
		return printJSON(client.CreateJob(context.Background(), positionals[1], body))
	case "job":
		if len(positionals) < 2 {
			fmt.Fprintln(os.Stderr, "job requires a job id")
			return 1
		}
		return printJSON(client.GetJob(context.Background(), positionals[1]))
	default:
		fmt.Fprintf(os.Stderr, "unknown command %s\n", command)
		return 1
	}
}

func parseArgs(argv []string) (map[string]string, []string, error) {
	flags := map[string]string{}
	positionals := make([]string, 0, len(argv))
	for index := 0; index < len(argv); index += 1 {
		token := argv[index]
		if token == "--base-url" || token == "--api-key" || token == "--cursor" || token == "--limit" {
			if index+1 >= len(argv) {
				return nil, nil, fmt.Errorf("missing value for %s", token)
			}
			flags[token[2:]] = argv[index+1]
			index += 1
			continue
		}
		if strings.HasPrefix(token, "-") {
			return nil, nil, fmt.Errorf("unknown flag %s", token)
		}
		positionals = append(positionals, token)
	}
	return flags, positionals, nil
}

func printJSON(value interface{}, err error) int {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if encodeErr := encoder.Encode(value); encodeErr != nil {
		fmt.Fprintln(os.Stderr, encodeErr.Error())
		return 1
	}
	return 0
}
