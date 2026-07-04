package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fuzzing-api/api"
	"fuzzing-api/load"
)

func main() {
	configPath := flag.String("config", "config/load.json", "load test configuration file")
	artifactsDir := flag.String("artifacts", "artifacts/load", "directory for load test artifacts")
	users := flag.Int("users", 2, "number of concurrent users")
	durationText := flag.String("duration", "", "run duration, for example 30s or 2m")
	requests := flag.Int64("requests", 10, "total requests to execute; use 0 with -duration for duration-only runs")
	scenarioFilter := flag.String("scenario", "", "comma-separated scenario names to run; empty runs all")
	flag.Parse()

	duration := time.Duration(0)
	if strings.TrimSpace(*durationText) != "" {
		parsed, err := time.ParseDuration(*durationText)
		if err != nil {
			exitf("invalid duration: %v", err)
		}
		duration = parsed
	}

	config, err := load.LoadConfig(*configPath)
	if err != nil {
		exitf("load config: %v", err)
	}

	scenarios := filterScenarios(config.Scenarios, *scenarioFilter)
	if len(scenarios) == 0 {
		exitf("no scenarios matched %q", *scenarioFilter)
	}

	if err := os.MkdirAll(*artifactsDir, 0755); err != nil {
		exitf("create artifacts directory: %v", err)
	}

	client := api.NewAPIClient(config.BaseURL)
	report, err := load.Run(context.Background(), client, scenarios, load.RunnerOptions{
		Users:    *users,
		Duration: duration,
		Requests: *requests,
	})
	if err != nil {
		exitf("run load test: %v", err)
	}

	if err := writeReport(*artifactsDir, report); err != nil {
		exitf("write report: %v", err)
	}

	fmt.Printf("Load artifacts saved in %s\n", *artifactsDir)
	fmt.Printf("Requests: %d, success: %d, failed: %d, rps: %.2f, avg: %.2fms, p95: %dms, p99: %dms\n",
		report.TotalRequests,
		report.Successful,
		report.Failed,
		report.RequestsPerSec,
		report.AvgMS,
		report.P95MS,
		report.P99MS,
	)

	if report.Failed > 0 {
		os.Exit(1)
	}
}

func filterScenarios(scenarios []load.Scenario, filter string) []load.Scenario {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return scenarios
	}

	allowed := map[string]bool{}
	for _, name := range strings.Split(filter, ",") {
		allowed[strings.TrimSpace(name)] = true
	}

	var selected []load.Scenario
	for _, scenario := range scenarios {
		if allowed[scenario.Name] {
			selected = append(selected, scenario)
		}
	}
	return selected
}

func writeReport(dir string, report *load.Report) error {
	summaryPath := filepath.Join(dir, "summary.json")
	samplesPath := filepath.Join(dir, "load-requests.jsonl")

	summaryFile, err := os.Create(summaryPath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(summaryFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		summaryFile.Close()
		return err
	}
	if err := summaryFile.Close(); err != nil {
		return err
	}

	samplesFile, err := os.Create(samplesPath)
	if err != nil {
		return err
	}
	encoder = json.NewEncoder(samplesFile)
	for _, sample := range report.Samples {
		if err := encoder.Encode(sample); err != nil {
			samplesFile.Close()
			return err
		}
	}
	return samplesFile.Close()
}

func exitf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
