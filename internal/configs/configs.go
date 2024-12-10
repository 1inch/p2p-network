// Package configs provides configuration utilities.
package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig is a generic function to load YAML configuration into any struct.
func LoadConfig[T any](filePath string) (*T, error) {
	var result *T
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close config file: %w", closeErr)
		}
	}()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	result = new(T)
	if decodeErr := decoder.Decode(result); decodeErr != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", decodeErr)
	}

	return result, err
}
