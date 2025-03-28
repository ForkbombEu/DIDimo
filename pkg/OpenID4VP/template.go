package OpenID4VP

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Variant represents the extracted format components.
type Variant struct {
	CredentialFormat string `json:"credential_format"`
	ClientIDScheme   string `json:"client_id_scheme"`
	RequestMethod    string `json:"request_method"`
	ResponseMode     string `json:"response_mode"`
}

type Variants struct {
	Variants []string `json:"variants"`
}

// Config represents the structure of the configuration file.
type Config struct {
	VariantKeys    map[string][]string `json:"variant_keys"`
	OptionalFields map[string]struct {
		Values   map[string][]string `json:"values"`
		Template string              `json:"template"`
	} `json:"optional_fields"`
}

// FinalFormat represents the full output structure.
type FinalFormat struct {
	Variant Variant     `json:"variant"`
	Form    interface{} `json:"form"`
}

// LoadJSON reads a JSON file and unmarshals it into a given struct.
func LoadJSON(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}
	return nil
}

// ValidateVariant checks if each value in the variant is within allowed keys.
func ValidateVariant(variant Variant, config Config) error {
	mapping := map[string]string{
		"credential_format": variant.CredentialFormat,
		"client_id_scheme":  variant.ClientIDScheme,
		"request_method":    variant.RequestMethod,
		"response_mode":     variant.ResponseMode,
	}

	for key, value := range mapping {
		if allowedValues, exists := config.VariantKeys[key]; exists {
			valid := false
			for _, allowed := range allowedValues {
				if value == allowed {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid value '%s' for key '%s'", value, key)
			}
		}
	}
	return nil
}

// ParseInput processes the input string and applies default and configuration JSON.
func ParseInput(input, defaultFile, configFile string) (*FinalFormat, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid input format")
	}

	// Create Variant struct
	variant := Variant{
		CredentialFormat: parts[0],
		ClientIDScheme:   parts[1],
		RequestMethod:    parts[2],
		ResponseMode:     parts[3],
	}

	// Load default JSON structure
	var defaultData map[string]any
	if err := LoadJSON(defaultFile, &defaultData); err != nil {
		return nil, err
	}

	// Load configuration JSON
	var config Config
	if err := LoadJSON(configFile, &config); err != nil {
		return nil, err
	}

	// Validate variant values
	if err := ValidateVariant(variant, config); err != nil {
		return nil, err
	}

	// Extract the "client" part of the default JSON
	client, ok := defaultData["form"].(map[string]any)["client"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid default JSON structure")
	}

	// Apply optional fields based on configuration
	variantMap := map[string]string{
		"credential_format": variant.CredentialFormat,
		"client_id_scheme":  variant.ClientIDScheme,
		"request_method":    variant.RequestMethod,
		"response_mode":     variant.ResponseMode,
	}

	for field, rule := range config.OptionalFields {
		for param, allowedValues := range rule.Values {
			if value, exists := variantMap[param]; exists {
				for _, allowed := range allowedValues {
					if value == allowed {
						client[field] = rule.Template
						break
					}
				}
			}
		}
	}

	// Update the form in the final structure
	defaultData["form"].(map[string]interface{})["client"] = client

	return &FinalFormat{
		Variant: variant,
		Form:    defaultData["form"],
	}, nil
}
