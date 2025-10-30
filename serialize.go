package dbml

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// ToJSON converts a Project to JSON bytes.
func (p *Project) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// FromJSON populates a Project from JSON bytes.
func (p *Project) FromJSON(data []byte) error {
	return json.Unmarshal(data, p)
}

// ToYAML converts a Project to YAML bytes.
func (p *Project) ToYAML() ([]byte, error) {
	return yaml.Marshal(p)
}

// FromYAML populates a Project from YAML bytes.
func (p *Project) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, p)
}
