package stepci

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type HTTP struct {
	Method   string            `yaml:"method"`
	URL      string            `yaml:"url"`
	Params   map[string]string `yaml:"params,omitempty"`
	Auth     any               `yaml:"auth,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty"`
	JSON     any               `yaml:"json,omitempty"`
	Captures map[string]struct {
		JSONPath string `yaml:"jsonpath"`
	} `yaml:"captures,omitempty"`
	Check struct {
		Status int `yaml:"status,omitempty"`
		Schema any `yaml:"schema,omitempty"`
	} `yaml:"check"`
}

type Step struct {
	Name string `yaml:"name"`
	HTTP HTTP   `yaml:"http"`
}

type Test struct {
	Steps []Step `yaml:"steps"`
}

type Config struct {
	Version    string          `yaml:"version"`
	Components map[string]any  `yaml:"components"`
	Tests      map[string]Test `yaml:"tests"`
}

// ConvertJSONToMap reads a JSON file and converts its contents to a map[string]any
func ConvertJSONToMap(jsonFilePath string) (map[string]any, error) {
	// Read the JSON file
	jsonFile, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	// Unmarshal JSON into a generic map
	var jsonData map[string]any
	if err := json.Unmarshal(jsonFile, &jsonData); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// Return the map representing the YAML structure
	return jsonData, nil
}

func GenerateYAML(config Config, filePath string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
