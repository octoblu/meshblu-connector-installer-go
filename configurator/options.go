package configurator

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/codegangsta/cli"
)

// OptionsConfig defines the service configurations
type OptionsConfig struct {
	Connector       string
	DownloadURI     string
	OutputDirectory string
	Legacy          bool
	ServiceName     string
	Hostname        string
	Port            int
	UUID, Token     string
}

// Options defines the service configurations
type Options interface {
	GetConnectorDirectory() string
	GetBinDirectory() string
	GetLogDirectory() string
	GetServiceName() string
	GetExecutablePath() string
	GetConnector() string
	GetDisplayName() string
	GetDescription() string
	GetPathEnv() string
	GetDownloadURI() string
	GetOutputDirectory() string
	GetLegacy() bool
	GetLegacyFlag() string
	GetHostname() string
	GetPort() int
	GetUUID() string
	GetToken() string
}

// NewOptions should create an options points
func NewOptions(context *cli.Context) Options {
	outputDirectory := context.String("output")
	if outputDirectory == "" {
		outputDirectory = GetDefaultServiceDirectory()
	}
	outputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		log.Fatalln("Invalid output directory:", err.Error())
	}
	return &OptionsConfig{
		Connector:       context.String("connector"),
		DownloadURI:     context.String("download-uri"),
		OutputDirectory: outputDirectory,
		Legacy:          context.Bool("legacy"),
		Hostname:        "meshblu.octoblu.com",
		Port:            443,
		UUID:            context.String("uuid"),
		Token:           context.String("token"),
	}
}

// GetConnectorDirectory gets the OS specific connector path
func (opts *OptionsConfig) GetConnectorDirectory() string {
	return path.Join(opts.OutputDirectory, opts.UUID)
}

// GetLogDirectory gets the OS specific log directory
func (opts *OptionsConfig) GetLogDirectory() string {
	return path.Join(opts.GetConnectorDirectory(), "log")
}

// GetBinDirectory gets the OS specific log directory
func (opts *OptionsConfig) GetBinDirectory() string {
	return path.Join(opts.OutputDirectory, "bin")
}

// GetConnector get connector name
func (opts *OptionsConfig) GetConnector() string {
	return opts.Connector
}

// GetDisplayName get service display name
func (opts *OptionsConfig) GetDisplayName() string {
	return fmt.Sprintf("MeshbluConnector %s", opts.GetUUID())
}

// GetDescription get service description
func (opts *OptionsConfig) GetDescription() string {
	return fmt.Sprintf("MeshbluConnector (%s) %s", opts.GetConnector(), opts.GetUUID())
}

// GetDownloadURI get download uri
func (opts *OptionsConfig) GetDownloadURI() string {
	return opts.DownloadURI
}

// GetOutputDirectory get output directory
func (opts *OptionsConfig) GetOutputDirectory() string {
	return opts.OutputDirectory
}

// GetLegacy get legacy bool
func (opts *OptionsConfig) GetLegacy() bool {
	return opts.Legacy
}

// GetLegacyFlag get legacy flag
func (opts *OptionsConfig) GetLegacyFlag() string {
	if opts.GetLegacy() {
		return "--legacy"
	}
	return ""
}

// GetHostname get meshblu hostname
func (opts *OptionsConfig) GetHostname() string {
	return opts.Hostname
}

// GetPort get meshblu port
func (opts *OptionsConfig) GetPort() int {
	return opts.Port
}

// GetUUID get meshblu uuid
func (opts *OptionsConfig) GetUUID() string {
	return opts.UUID
}

// GetToken get meshblu token
func (opts *OptionsConfig) GetToken() string {
	return opts.Token
}