package load

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BaseURL   string     `json:"baseURL"`
	Scenarios []Scenario `json:"scenarios"`
}

type Scenario struct {
	Name     string          `json:"name"`
	Method   string          `json:"method"`
	Endpoint string          `json:"endpoint"`
	Body     json.RawMessage `json:"body,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.BaseURL) == "" {
		return fmt.Errorf("baseURL is required")
	}
	if len(c.Scenarios) == 0 {
		return fmt.Errorf("at least one scenario is required")
	}

	for i, scenario := range c.Scenarios {
		if strings.TrimSpace(scenario.Name) == "" {
			return fmt.Errorf("scenario %d name is required", i)
		}
		if strings.TrimSpace(scenario.Method) == "" {
			return fmt.Errorf("scenario %q method is required", scenario.Name)
		}
		if strings.TrimSpace(scenario.Endpoint) == "" {
			return fmt.Errorf("scenario %q endpoint is required", scenario.Name)
		}
	}

	return nil
}
