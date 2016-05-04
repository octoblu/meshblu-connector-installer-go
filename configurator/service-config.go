package configurator

import (
	"encoding/json"
	"path/filepath"
)

// ServiceConfig is the runner connector config structure.
type ServiceConfig struct {
	ServiceName   string
	DisplayName   string
	Description   string
	ConnectorName string
	BinPath       string
	Legacy        bool
	Dir           string

	Stderr, Stdout string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(opts Options) *ServiceConfig {
	return &ServiceConfig{
		ServiceName:   opts.GetServiceName(),
		DisplayName:   opts.GetDisplayName(),
		Description:   opts.GetDescription(),
		ConnectorName: opts.GetConnector(),
		BinPath:       opts.GetBinDirectory(),
		Legacy:        opts.GetLegacy(),
		Dir:           opts.GetConnectorDirectory(),
		Stderr:        filepath.Join(opts.GetLogDirectory(), "connector-error.log"),
		Stdout:        filepath.Join(opts.GetLogDirectory(), "connector.log"),
	}
}

// ToJSON serializes the object to the meshblu.json format
func (config *ServiceConfig) ToJSON() ([]byte, error) {
	return json.Marshal(config)
}
