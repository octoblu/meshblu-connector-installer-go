package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/meshblu-connector-installer-go/configurator"
	"github.com/octoblu/meshblu-connector-installer-go/downloader"
	"github.com/octoblu/meshblu-connector-installer-go/extractor"
	"github.com/octoblu/meshblu-connector-installer-go/foreverizer"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-installer:main")

func main() {
	app := cli.NewApp()
	app.Name = "meshblu-connector-installer"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "connector, c",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_CONNECTOR",
			Usage:  "Meshblu connector name",
		},
		cli.StringFlag{
			Name:   "hostname, -h",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_HOSTNAME",
			Usage:  "Meshblu device hostname",
		},
		cli.StringFlag{
			Name:   "output, o",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_OUTPUT",
			Usage:  "Output directory",
		},
		cli.IntFlag{
			Name:   "port",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_PORT",
			Usage:  "Meshblu device port",
		},
		cli.StringFlag{
			Name:   "uuid, -u",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_UUID",
			Usage:  "Meshblu device uuid",
		},
		cli.StringFlag{
			Name:   "tag, t",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_TAG",
			Usage:  "Tag version. Defaults to 'latest'",
		},
		cli.StringFlag{
			Name:   "token",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_TOKEN",
			Usage:  "Meshblu device token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	connector, hostname, outputDirectory, port, uuid, tag, token := getOpts(context)
	platform := "osx"
	err := os.MkdirAll(outputDirectory, 0755)
	fatalIfError("Error creating output directory", err)

	baseURI := "https://meshblu-connector.octoblu.com"
	downloadClient := downloader.New(outputDirectory, baseURI)
	downloadFile, err := downloadClient.DownloadConnector(connector, tag, platform)
	fatalIfError("Error downloading", err)

	extractorClient := extractor.New()
	err = extractorClient.Do(downloadFile, outputDirectory)
	fatalIfError("Error extracting:", err)

	configuratorClient := configurator.New(outputDirectory)
	err = configuratorClient.WriteMeshblu(uuid, token, hostname, port)
	fatalIfError("Error writing meshblu config:", err)

	foreverizerClient := foreverizer.New()
	err = foreverizerClient.Do(uuid, outputDirectory)
	fatalIfError("Error setuping device to run forever", err)
}

func getOpts(context *cli.Context) (string, string, string, int, string, string, string) {
	connector := context.String("connector")
	hostname := context.String("hostname")
	output := context.String("output")
	port := context.Int("port")
	uuid := context.String("uuid")
	tag := context.String("tag")
	token := context.String("token")

	if connector == "" || uuid == "" || token == "" {
		cli.ShowAppHelp(context)

		if connector == "" {
			color.Red("  Missing required flag --connector or MESHBLU_CONNECTOR_INSTALLER_CONNECTOR")
		}

		if uuid == "" {
			color.Red("  Missing required flag --uuid or MESHBLU_CONNECTOR_INSTALLER_OUTPUT")
		}

		if token == "" {
			color.Red("  Missing required flag --token or MESHBLU_CONNECTOR_INSTALLER_OUTPUT")
		}
		os.Exit(1)
	}

	if output == "" {
		output = path.Join(os.Getenv("HOME"), "Library", "Application Support", "Octoblu", uuid)
	}

	outputDirectory, err := filepath.Abs(output)
	if err != nil {
		log.Fatalln("Invalid output directory:", err.Error())
	}

	if hostname == "" {
		hostname = "meshblu.octoblu.com"
	}

	if port == 0 {
		port = 443
	}

	if tag == "" {
		tag = "latest"
	}

	return connector, hostname, outputDirectory, port, uuid, tag, token
}

func fatalIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Fatalln(msg, err.Error())
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
