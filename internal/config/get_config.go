package config

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

func GetConfig() (*ConfigStruct, error) {
	configPath := os.Getenv("CONFIG_YAML_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("environment variable CONFIG_PATH is not set")
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	expanded, err := expandEnvVars(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to expand env vars in config: %w", err)
	}

	var cfg ConfigStruct
	if err := yaml.Unmarshal(expanded, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func expandEnvVars(raw []byte) ([]byte, error) {
	tmpl, err := template.New("config").Option("missingkey=error").Parse(string(raw))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, envMap()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func envMap() map[string]string {
	env := map[string]string{}
	for _, pair := range os.Environ() {
		varParts := bytes.SplitN([]byte(pair), []byte("="), 2)
		if len(varParts) == 2 {
			env[string(varParts[0])] = string(varParts[1])
		}
	}
	return env
}