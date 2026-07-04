package load

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"fuzzing-api/api"
)

type RunnerOptions struct {
	Users    int
	Duration time.Duration
	Requests int64
}

type Sample struct {
	Timestamp  string `json:"timestamp"`
	Scenario   string `json:"scenario"`
	Method     string `json:"method"`
	Endpoint   string `json:"endpoint"`
	StatusCode int    `json:"statusCode"`
	DurationMS int64  `json:"durationMs"`
	Error      string `json:"error,omitempty"`
}

type Report struct {
	StartedAt       string         `json:"startedAt"`
	FinishedAt      string         `json:"finishedAt"`
	Users           int            `json:"users"`
	DurationSeconds float64        `json:"durationSeconds"`
	TotalRequests   int            `json:"totalRequests"`
	Successful      int            `json:"successful"`
	Failed          int            `json:"failed"`
	RequestsPerSec  float64        `json:"requestsPerSec"`
	MinMS           int64          `json:"minMs"`
	AvgMS           float64        `json:"avgMs"`
	P95MS           int64          `json:"p95Ms"`
	P99MS           int64          `json:"p99Ms"`
	MaxMS           int64          `json:"maxMs"`
	StatusCodes     map[string]int `json:"statusCodes"`
	Samples         []Sample       `json:"-"`
}

func Run(ctx context.Context, client *api.APIClient, scenarios []Scenario, options RunnerOptions) (*Report, error) {
	if client == nil {
		return nil, fmt.Errorf("client is required")
	}
	if len(scenarios) == 0 {
		return nil, fmt.Errorf("at least one scenario is required")
	}
	if options.Users <= 0 {
		options.Users = 1
	}
	if options.Duration <= 0 && options.Requests <= 0 {
		options.Requests = int64(len(scenarios))
	}

	runCtx := ctx
	cancel := func() {}
	if options.Duration > 0 {
		runCtx, cancel = context.WithTimeout(ctx, options.Duration)
	}
	defer cancel()

	started := time.Now()
	var issued int64
	samples := make(chan Sample, options.Users*2)
	var wg sync.WaitGroup

	for user := 0; user < options.Users; user++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			for {
				current := atomic.AddInt64(&issued, 1)
				if options.Requests > 0 && current > options.Requests {
					return
				}

				select {
				case <-runCtx.Done():
					return
				default:
				}

				scenario := scenarios[(int(current)-1+offset)%len(scenarios)]
				samples <- executeScenario(client, scenario)
			}
		}(user)
	}

	go func() {
		wg.Wait()
		close(samples)
	}()

	report := &Report{
		StartedAt:   started.Format(time.RFC3339),
		Users:       options.Users,
		StatusCodes: map[string]int{},
	}
	var durations []int64
	var totalDuration int64

	for sample := range samples {
		report.Samples = append(report.Samples, sample)
		report.TotalRequests++
		durations = append(durations, sample.DurationMS)
		totalDuration += sample.DurationMS

		if sample.Error == "" && sample.StatusCode >= 200 && sample.StatusCode < 400 {
			report.Successful++
		} else {
			report.Failed++
		}
		report.StatusCodes[statusKey(sample)]++
	}

	finished := time.Now()
	report.FinishedAt = finished.Format(time.RFC3339)
	report.DurationSeconds = finished.Sub(started).Seconds()
	if report.DurationSeconds > 0 {
		report.RequestsPerSec = float64(report.TotalRequests) / report.DurationSeconds
	}
	if report.TotalRequests > 0 {
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		report.MinMS = durations[0]
		report.MaxMS = durations[len(durations)-1]
		report.AvgMS = float64(totalDuration) / float64(report.TotalRequests)
		report.P95MS = percentile(durations, 0.95)
		report.P99MS = percentile(durations, 0.99)
	}

	return report, nil
}

func executeScenario(client *api.APIClient, scenario Scenario) Sample {
	method := strings.ToUpper(strings.TrimSpace(scenario.Method))
	body := ""
	if len(scenario.Body) > 0 && string(scenario.Body) != "null" {
		var compact bytes.Buffer
		if err := json.Compact(&compact, scenario.Body); err == nil {
			body = compact.String()
		} else {
			body = string(scenario.Body)
		}
	}

	resp, statusCode, duration, err := client.Request(method, scenario.Endpoint, body)
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	sample := Sample{
		Timestamp:  time.Now().Format(time.RFC3339),
		Scenario:   scenario.Name,
		Method:     method,
		Endpoint:   scenario.Endpoint,
		StatusCode: statusCode,
		DurationMS: duration.Milliseconds(),
	}
	if err != nil {
		sample.Error = err.Error()
	}
	return sample
}

func percentile(values []int64, rank float64) int64 {
	if len(values) == 0 {
		return 0
	}
	index := int(rank*float64(len(values)) + 0.5)
	if index < 1 {
		index = 1
	}
	if index > len(values) {
		index = len(values)
	}
	return values[index-1]
}

func statusKey(sample Sample) string {
	if sample.Error != "" && sample.StatusCode == 0 {
		return "error"
	}
	return fmt.Sprintf("%d", sample.StatusCode)
}
